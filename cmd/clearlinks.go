package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jonesrussell/page-prowler/redis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clearlinksCmd = &cobra.Command{
	Use:   "clearlinks",
	Short: "Clear the Redis set for a given crawlsiteid",
	Run: func(cmd *cobra.Command, args []string) {
		crawlsiteid := viper.GetString("crawlsiteid")
		if crawlsiteid == "" {
			fmt.Println("crawlsiteid is required")
			os.Exit(1)
		}

		ctx := context.Background()
		redisClient, err := redis.NewClient(viper.GetString("REDIS_HOST"), viper.GetString("REDIS_AUTH"), viper.GetString("REDIS_PORT"))
		if err != nil {
			fmt.Println("Failed to initialize Redis client", "error", err)
			os.Exit(1)
		}

		_, err = redisClient.Del(ctx, crawlsiteid)
		if err != nil {
			fmt.Println("Failed to clear Redis set", "error", err)
			os.Exit(1)
		}

		fmt.Println("Redis set cleared successfully")
	},
}

func init() {
	rootCmd.AddCommand(clearlinksCmd)
}
