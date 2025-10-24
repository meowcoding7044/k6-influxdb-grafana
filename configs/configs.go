package configs

import (
	"context"
	"fmt"
	"goFirst1/logs"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ------------------------------- //
// Viper Config                    //
// ------------------------------- //
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func InitTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

// ------------------------------- //
// Database Config                    //
// ------------------------------- //
type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func InitDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Bangkok",
		viper.GetString("db.host"),
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.database"),
		viper.GetInt("db.port"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.Error("Database connection failed: " + err.Error())
		panic(err)
	}

	// err = db.AutoMigrate()
	// if err != nil {
	// 	panic(err)
	// }

	Database = DbInstance{Db: db}
	logs.Info("✅ Database connected successfully")
	// defer db.Close()
	return db
}

// ------------------------------- //
// Redis Config                    //
// ------------------------------- //
type RedisInstance struct {
	Client *redis.Client
}

var Redis RedisInstance

func InitRedis() *redis.Client {
	addr := fmt.Sprintf("%s:%d",
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
	)

	client := redis.NewClient(&redis.Options{
		Addr:             addr,
		Password:         viper.GetString("redis.password"),
		DB:               viper.GetInt("redis.db"),
		DisableIndentity: true, //ปิด identity handshake (เพื่อข้าม maint_notifications)
		// ปิด internal logger
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			return nil
		},
	})

	//auto reconnect + Retry Backoff
	ctx := context.Background()
	maxRetries := 5
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := client.Ping(ctx).Err()
		if err == nil {
			Redis = RedisInstance{Client: client}
			logs.Info(fmt.Sprintf("✅ Redis connected successfully at %s", addr))
			return client
		}

		logs.Error(fmt.Sprintf("⚠️ Redis connection failed (attempt %d/%d): %v", attempt, maxRetries, err))
		time.Sleep(retryDelay)
		retryDelay *= 2 // Exponential backoff
	}

	logs.Error("❌ Redis connection failed after retries — shutting down.")
	panic("Redis initialization failed")
}
