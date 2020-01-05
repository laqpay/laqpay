package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"

	"github.com/laqpay/laqpay/src/cipher"
	"github.com/laqpay/laqpay/src/util/apputil"
	"github.com/laqpay/laqpay/src/visor"
	"github.com/laqpay/laqpay/src/visor/dbutil"
)

const (
	blockchainPubkey = "022696ad3d8bd010e6cfdb7fede6cdf5b63a7b0ba5f861abce99642f7c8bf2cac0"
)

// wrapDB calls dbutil.WrapDB and disables all logging
func wrapDB(db *bolt.DB) *dbutil.DB {
	wdb := dbutil.WrapDB(db)
	wdb.ViewLog = false
	wdb.ViewTrace = false
	wdb.UpdateLog = false
	wdb.UpdateTrace = false
	wdb.DurationLog = false
	return wdb
}

func checkDBCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Verify the database",
		Use:   "checkdb [db path]",
		Long: `Checks if the given database file contains valid laqpay blockchain data.
    If no argument is specificed, the default data.db in $HOME/.$COIN/ will be checked.`,
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  checkDB,
	}
}

func checkDB(_ *cobra.Command, args []string) error {
	// get db path
	dbPath := ""
	if len(args) > 0 {
		dbPath = args[0]
	}
	dbPath, err := resolveDBPath(cliConfig, dbPath)
	if err != nil {
		return err
	}

	// check if this file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("db file: %v does not exist", dbPath)
	}

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout:  5 * time.Second,
		ReadOnly: true,
	})

	if err != nil {
		return fmt.Errorf("open db failed: %v", err)
	}
	pubkey, err := cipher.PubKeyFromHex(blockchainPubkey)
	if err != nil {
		return fmt.Errorf("decode blockchain pubkey failed: %v", err)
	}

	go func() {
		apputil.CatchInterrupt(quitChan)
	}()

	if err := visor.CheckDatabase(wrapDB(db), pubkey, quitChan); err != nil {
		if err == visor.ErrVerifyStopped {
			return nil
		}
		return fmt.Errorf("checkdb failed: %v", err)
	}

	fmt.Println("check db success")
	return nil
}

func checkDBEncodingCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Verify the database data encoding",
		Use:   "checkDBDecoding [db path]",
		Long: `Verify the generated binary encoders match the dynamic encoders for database data.
    If no argument is specificed, the default data.db in $HOME/.$COIN/ will be checked.`,
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  checkDBDecoding,
	}
}

func checkDBDecoding(_ *cobra.Command, args []string) error {
	// get db path
	dbPath := ""
	if len(args) > 0 {
		dbPath = args[0]
	}
	dbPath, err := resolveDBPath(cliConfig, dbPath)
	if err != nil {
		return err
	}

	// check if this file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("db file: %v does not exist", dbPath)
	}

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout:  5 * time.Second,
		ReadOnly: true,
	})

	if err != nil {
		return fmt.Errorf("open db failed: %v", err)
	}

	go func() {
		apputil.CatchInterrupt(quitChan)
	}()

	if err := visor.VerifyDBLaqencoderSafe(wrapDB(db), quitChan); err != nil {
		if err == visor.ErrVerifyStopped {
			return nil
		}
		return fmt.Errorf("checkDBDecoding failed: %v", err)
	}

	fmt.Println("check db decoding success")
	return nil

}
