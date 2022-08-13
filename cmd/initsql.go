/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
type DB struct {
	DbName    string `mapstructure:"db_name"`
	TableFile string `mapstructure:"table_file"`
	InitFile  string `mapstructure:"init_file"`
}
type SqlInfo struct {
	DbConfig DbConfig      `mapstructure:"db_server"`
	InitDb   map[string]DB `mapstructure:"init_db"`
}

func initConfig() *SqlInfo {
	v := viper.New()
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	SqlInfo := SqlInfo{}
	if err := v.Unmarshal(&SqlInfo); err != nil {
		panic(err)
	}
	return &SqlInfo
}

// initsqlCmd represents the initsql command
var initsqlCmd = &cobra.Command{
	Use:   "initsql",
	Short: "A brief description of your command",
	Long:  `A brief description of your command`,
	Run: func(cmd *cobra.Command, args []string) {
		SqlInfo := initConfig()
		host := SqlInfo.DbConfig.Host
		port := SqlInfo.DbConfig.Port
		username := SqlInfo.DbConfig.Username
		password := SqlInfo.DbConfig.Password
		//dsn := "root:@tcp(localhost:3306)/sys?charset=utf8&parseTime=True"
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, strconv.Itoa(port))
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("DSN : %s Format failed, Error: %v \n", dsn, err)
			panic(err)
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			log.Printf("Connection %s Failed, Error: %v \n", dsn, err)
			return
		}
		log.Println("数据库连接成功")
		for _, v := range SqlInfo.InitDb {
			createdb := fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8 COLLATE utf8_general_ci;", v.DbName)
			//rows, err := db.Exec(createdb)
			_, err := db.Exec(createdb)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initsqlCmd)
}
