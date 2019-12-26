package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/rtravitz/chatter/database"
	"github.com/rtravitz/chatter/schema"
	"github.com/rtravitz/chatter/user"
	"log"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:0.0.0.0"`
			Name       string `conf:"default:chatter"`
			DisableTLS bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "CHATTER", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return err
			}
			fmt.Println(usage)
			return nil
		}
		return err
	}

	dbConfig := database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	var err error
	switch cfg.Args.Num(0) {
	case "migrate":
		err = migrate(dbConfig)
	case "useradd":
		err = useradd(dbConfig, cfg.Args.Num(1), cfg.Args.Num(2))
	default:
		err = errors.New("Must specify a command")
	}

	if err != nil {
		return err
	}

	return nil
}

func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return err
	}

	fmt.Println("Migrations complete")
	return nil
}

func useradd(cfg database.Config, email, password string) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}

	defer db.Close()

	if email == "" || password == "" {
		return errors.New("useradd must be called with arguments for username and password")
	}

	fmt.Printf("user will be created with email %q and password %q\n", email, password)

	fmt.Print("Continue? (1/0) ")
	var confirm bool
	if _, err := fmt.Scanf("%t\n", &confirm); err != nil {
		return err
	}

	if !confirm {
		fmt.Println("Canceling")
		return nil
	}

	nu := user.NewUser{
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
	}

	u, err := user.Create(db, nu, time.Now())
	if err != nil {
		return err
	}

	fmt.Println("User created with id:", u.ID)

	return nil
}
