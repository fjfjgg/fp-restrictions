package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"encoding/json"
	"reflect"

	"github.com/gofrp/fp-restrictions/pkg/server"
	"github.com/gofrp/fp-restrictions/pkg/server/controller"

	"github.com/spf13/cobra"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
)

const version = "0.1.1"

var (
	showVersion       bool
	showRestrictions  bool
	bindAddr          string
	restrictionFile   string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version")
	rootCmd.PersistentFlags().BoolVarP(&showRestrictions, "show_restrictions", "s", false, "show restrictions")
	rootCmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "l", "127.0.0.1:7200", "bind address")
	rootCmd.PersistentFlags().StringVarP(&restrictionFile, "restriction_file", "f", "./restrictions", "restriction file")
}

var rootCmd = &cobra.Command{
	Use:   "fp-restrictions",
	Short: "fp-restrictions is the server plugin of frp to support restrictions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version)
			return nil
		}
		if showRestrictions {
			test := &plugin.NewProxyContent{}
			t := reflect.Indirect(reflect.ValueOf(test.NewProxy)).Type()
			for i := 0; i < t.NumField(); i++ {
		 		f := t.Field(i)
				fmt.Println(f.Name, f.Type)
			}
			return nil
		}
		restrictions, err := ParseRestrictionsFromFile(restrictionFile)
		if err != nil {
			log.Printf("parse tokens from file %s error: %v", restrictionFile, err)
			return err
		}
		s, err := server.New(server.Config{
			BindAddress:  bindAddr,
			Restrictions: restrictions,
		})
		if err != nil {
			return err
		}
		s.Run()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ParseRestrictionsFromFile(file string) (map[string]controller.Restriction, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]controller.Restriction)
	rows := strings.Split(string(buf), "\n")
	for _, row := range rows {
		kvs := strings.SplitN(row, "=", 2)
		if len(kvs) == 2 {
			user := strings.TrimSpace(kvs[0])
			tr := strings.SplitN(kvs[1],":",2)
			var r controller.Restriction
			r.Token = strings.TrimSpace(tr[0])
			if len(tr) == 2 {
				if err := json.Unmarshal([]byte(tr[1]), &r.Restriction); err != nil {
					log.Println(err)
				}
			}
			ret[user] = r
		}
	}
	return ret, nil
}
