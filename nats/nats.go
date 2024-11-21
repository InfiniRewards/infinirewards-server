package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"infinirewards/utils"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/jwt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nkeys"
)

var NC *nats.Conn
var js jetstream.JetStream

var accountPublicKey string
var accSeed string

type Account struct {
	Seed      string
	PublicKey string
}

func ConnectNats() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	accountPublicKey = os.Getenv("NATS_ACCOUNT_PUBLIC_KEY")
	accSeed = os.Getenv("NATS_ACCOUNT_SEED")
	natsUrl := os.Getenv("NATS_URL")

	// Validate required environment variables
	if accountPublicKey == "" || accSeed == "" || natsUrl == "" {
		return fmt.Errorf("missing required NATS environment variables: URL=%v, PublicKey=%v, Seed=%v",
			natsUrl != "", accountPublicKey != "", accSeed != "")
	}

	// Check if credentials file exists
	if _, err := os.Stat("./backend.creds"); err != nil {
		return fmt.Errorf("NATS credentials file not found at ./backend.creds: %w", err)
	}

	NC, err = nats.Connect(natsUrl,
		nats.UserCredentials("./backend.creds"),
		nats.Name("InfiniRewards Microservice Server (backend)"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second*5),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS server at %s: %w", natsUrl, err)
	}

	js, err = jetstream.New(NC)
	if err != nil {
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}

	_, err = js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:        "webhooks",
		Description: "Stream of webhook requests",
		Subjects:    []string{"webhooks.>"},
		MaxBytes:    -1,
		Retention:   jetstream.WorkQueuePolicy,
		AllowDirect: true,
		Replicas:    3,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update webhooks stream (replicas: 3): %w", err)
	}

	_, err = js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "token",
		Description: "JWTs",
		MaxBytes:    -1,
		TTL:         time.Hour * 24 * 90,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update token KV bucket: %w", err)
	}

	_, err = js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "users",
		Description: "Users",
		MaxBytes:    -1,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update users KV bucket: %w", err)
	}

	_, err = js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "merchants",
		Description: "Merchants",
		MaxBytes:    -1,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update merchants KV bucket: %w", err)
	}

	_, err = js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "apikeys",
		Description: "API keys",
		MaxBytes:    -1,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update API keys KV bucket: %w", err)
	}

	_, err = js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "phoneVerification",
		Description: "Phone verifications",
		MaxBytes:    -1,
		TTL:         time.Minute * 5,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update phone verification KV bucket: %w", err)
	}

	return nil
}

func PutKV(ctx context.Context, bucket string, key string, value []byte) error {
	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to get KV bucket: %w", err)
	}

	_, err = kv.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to put KV value: %w", err)
	}

	return nil
}

func GetKV(ctx context.Context, bucket string, key string) (jetstream.KeyValueEntry, error) {
	if js == nil {
		return nil, fmt.Errorf("JetStream not initialized")
	}

	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get KV bucket: %w", err)
	}

	data, err := kv.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get KV value: %w", err)
	}

	return data, nil
}

func RemoveKV(ctx context.Context, bucket string, key string) error {
	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to get KV bucket: %w", err)
	}

	err = kv.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete KV value: %w", err)
	}

	return nil
}

func GetKVValues[T any](ctx context.Context, bucket string, keyFilter string, transformFunc func(jetstream.KeyValueEntry, *T)) ([]*T, error) {
	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get KV bucket: %w", err)
	}

	watcher, _ := kv.Watch(ctx, keyFilter, jetstream.IgnoreDeletes())
	defer watcher.Stop()

	result := make([]*T, 0)

	for entry := range watcher.Updates() {
		if entry == nil {
			break
		}
		var t T
		if err := json.Unmarshal(entry.Value(), &t); err != nil {
			return nil, fmt.Errorf("failed to unmarshal KV value: %w", err)
		}
		transformFunc(entry, &t)
		result = append(result, &t)
	}
	return result, nil
}

func PutObject(ctx context.Context, bucket string, key string, value []byte) (*jetstream.ObjectInfo, error) {
	os, err := js.ObjectStore(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get object store: %w", err)
	}
	return os.PutBytes(ctx, key, value)

}

func GetObject(ctx context.Context, bucket string, key string) (jetstream.ObjectResult, error) {
	os, err := js.ObjectStore(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get object store: %w", err)
	}
	data, err := os.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return data, nil
}

func RemoveObject(ctx context.Context, bucket string, key string) error {
	os, err := js.ObjectStore(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to get object store: %w", err)
	}
	err = os.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

func PublishStream(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return js.PublishMsg(ctx, msg, opts...)
}

func GetStreamMsg(ctx context.Context, streamName string, subjectFilter string) (*jetstream.Msg, error) {
	randomness, err := utils.GenerateRandomString(5)
	if err != nil {
		return nil, err
	}
	consumer, err := js.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Name:           "lookup" + randomness,
		AckPolicy:      jetstream.AckNonePolicy,
		DeliverPolicy:  jetstream.DeliverLastPolicy,
		FilterSubjects: []string{subjectFilter},
	})
	if err != nil {
		return nil, err
	}
	msgs, err := consumer.FetchNoWait(1)
	if err != nil {
		return nil, err
	}
	var result *jetstream.Msg
	for msg := range msgs.Messages() {
		result = &msg
	}

	return result, nil
}

func LookupStream[T any](ctx context.Context, streamName string, subjectFilters []string, transformFunc func(jetstream.Msg, *T)) ([]*T, error) {
	randomness, err := utils.GenerateRandomString(5)
	if err != nil {
		return nil, err
	}
	consumer, err := js.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Name:           "lookup" + randomness,
		AckPolicy:      jetstream.AckNonePolicy,
		DeliverPolicy:  jetstream.DeliverLastPerSubjectPolicy,
		FilterSubjects: subjectFilters,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*T, 0)

	remaining := consumer.CachedInfo().NumPending
	iter, _ := consumer.Messages()
	for ; remaining > 0; remaining-- {
		msg, err := iter.Next()
		// Next can return error, e.g. when iterator is closed or no heartbeats were received
		if err != nil {
			return nil, err
		}
		var t T
		if err := json.Unmarshal(msg.Data(), &t); err != nil {
			return nil, err
		}
		transformFunc(msg, &t)
		result = append(result, &t)
	}
	iter.Stop()
	return result, nil
}

func DeleteStreamMsg(ctx context.Context, streamName string, seq uint64) error {
	stream, err := js.Stream(ctx, streamName)
	if err != nil {
		return err
	}
	stream.DeleteMsg(ctx, seq)
	return nil
}

func GenerateUserCredsString(ctx context.Context, service string, userPublicKey string, userSeed string) ([]byte, string, error) {
	accountSigningKey, err := nkeys.ParseDecoratedNKey([]byte(accSeed))
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse account signing key: %w", err)
	}

	userJWT, err := GenerateUserJWT(userPublicKey, accountPublicKey, accountSigningKey, service)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate user JWT: %w", err)
	}

	userSigningKey, err := nkeys.ParseDecoratedNKey([]byte(userSeed))
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse user signing key: %w", err)
	}

	userSeedBytes, err := userSigningKey.Seed()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user seed bytes: %w", err)
	}

	credsContent, err := jwt.FormatUserConfig(userJWT, userSeedBytes)
	if err != nil {
		return nil, "", fmt.Errorf("failed to format user config: %w", err)
	}

	return credsContent, userJWT, nil
}

func GenerateUserKey() (userPublicKey string, userSeed []byte) {
	kp, err := nkeys.CreateUser()
	if err != nil {
		return "", nil
	}
	if userSeed, err = kp.Seed(); err != nil {
		return "", nil
	} else if userPublicKey, err = kp.PublicKey(); err != nil {
		return "", nil
	}
	return
}

func GenerateUserJWT(userPublicKey, accountPublicKey string, accountSigningKey nkeys.KeyPair, service string) (userJWT string, err error) {
	uc := jwt.NewUserClaims(userPublicKey)
	switch service {
	case "auth":
		uc.Pub.Allow.Add(fmt.Sprintf("auth.*.%s.>", userPublicKey)) // auth
		uc.Pub.Allow.Add(fmt.Sprintf("auth.*.%s", userPublicKey))   // auth
	case "infinirewards":
		uc.Pub.Allow.Add(fmt.Sprintf("infinirewards.*.%s.>", userPublicKey))
		uc.Pub.Allow.Add(fmt.Sprintf("infinirewards.*.%s", userPublicKey))
	}
	// uc.Resp = &jwt.ResponsePermission{MaxMsgs: 1}
	uc.Sub.Allow.Add(fmt.Sprintf("_INBOX.%s.>", userPublicKey))
	uc.Sub.Allow.Add(fmt.Sprintf("_INBOX.%s", userPublicKey))
	uc.Expires = time.Now().Add(time.Hour).Unix() // expire in an hour
	uc.IssuerAccount = accountPublicKey
	vr := jwt.ValidationResults{}
	uc.Validate(&vr)
	if vr.IsBlocking(true) {
		err = fmt.Errorf("generated user claim is invalid")
		return
	}
	userJWT, err = uc.Encode(accountSigningKey)
	return
}

// CreateKVBuckets creates all necessary KV buckets for testing
func CreateKVBuckets() error {
	js, err := jetstream.New(NC)
	if err != nil {
		return err
	}

	buckets := []string{"apikeys", "phoneVerification", "token", "users"}
	for _, bucket := range buckets {
		_, err := js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
			Bucket: bucket,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// CleanupKVBuckets deletes all KV buckets used in testing
func CleanupKVBuckets() error {
	js, err := jetstream.New(NC)
	if err != nil {
		return err
	}

	buckets := []string{"apikeys", "phoneVerification", "token", "users"}
	for _, bucket := range buckets {
		kv, err := js.KeyValue(context.Background(), bucket)
		if err == nil {
			// Get all keys in the bucket
			keys, err := kv.Keys(context.Background())
			if err == nil {
				// Delete each key individually
				for _, key := range keys {
					kv.Delete(context.Background(), key)
				}
			}
			// Delete the bucket itself
			js.DeleteKeyValue(context.Background(), bucket)
		}
	}

	return nil
}

// Close closes the NATS connection
func Close() {
	if NC != nil {
		NC.Close()
	}
}

func ConnectNatsTest() error {
	var err error

	// Connect to NATS
	NC, err = nats.Connect("nats://127.0.0.1:4222",
		nats.Name("InfiniRewards Test Server"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second*5),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}

	// Initialize JetStream
	js, err = jetstream.New(NC)
	if err != nil {
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create streams and buckets
	if err := createTestStreams(); err != nil {
		return fmt.Errorf("failed to create streams: %w", err)
	}

	if err := createTestBuckets(); err != nil {
		return fmt.Errorf("failed to create buckets: %w", err)
	}

	return nil
}

// Add these helper functions
func createTestStreams() error {
	// Create webhooks stream
	_, err := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:        "webhooks",
		Description: "Stream of webhook requests",
		Subjects:    []string{"webhooks.>"},
		MaxBytes:    -1,
		Retention:   jetstream.WorkQueuePolicy,
		AllowDirect: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create webhooks stream: %w", err)
	}

	return nil
}

func createTestBuckets() error {
	// Create KV buckets
	buckets := []struct {
		name        string
		description string
		ttl         time.Duration
	}{
		{"token", "JWTs", time.Hour * 24 * 90},
		{"users", "Users", 0},
		{"merchants", "Merchants", 0},
		{"apikeys", "API keys", 0},
		{"phoneVerification", "Phone verifications", time.Minute * 5},
	}

	for _, bucket := range buckets {
		_, err := js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
			Bucket:      bucket.name,
			Description: bucket.description,
			MaxBytes:    -1,
			TTL:         bucket.ttl,
		})
		if err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", bucket.name, err)
		}
	}

	return nil
}
