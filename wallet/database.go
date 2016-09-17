// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wallet

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
)

// Wallet is a connection to a Factom Wallet Database
type Wallet struct {
	*WalletDatabaseOverlay
	transactions map[string]*factoid.Transaction
	txdb         *TXDatabaseOverlay
}

func (w *Wallet) InitWallet() error {
	dbSeed, err := w.GetDBSeed()
	if err != nil {
		return err
	}
	if dbSeed == nil {
		seed := make([]byte, 64)
		if n, err := rand.Read(seed); err != nil {
			return err
		} else if n != 64 {
			return fmt.Errorf("Wrong number of bytes read: %d", n)
		}
		err = w.InsertDBSeed(seed)
		if err != nil {
			return err
		}
		err = w.InsertNextDBSeed(seed)
		if err != nil {
			return err
		}
	}
	nextSeed, err := w.GetNextDBSeed()
	if err != nil {
		return err
	}
	if nextSeed == nil {
		return fmt.Errorf("Database does not contain nextSeed!")
	}
	return nil
}

func NewOrOpenLevelDBWallet(path string) (*Wallet, error) {
	w := new(Wallet)
	w.transactions = make(map[string]*factoid.Transaction)
	db, err := NewLevelDB(path)
	if err != nil {
		return nil, err
	}
	w.WalletDatabaseOverlay = db
	err = w.InitWallet()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func NewOrOpenBoltDBWallet(path string) (*Wallet, error) {
	w := new(Wallet)
	w.transactions = make(map[string]*factoid.Transaction)
	db, err := NewBoltDB(path)
	if err != nil {
		return nil, err
	}
	w.WalletDatabaseOverlay = db
	err = w.InitWallet()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func NewMapDBWallet() (*Wallet, error) {
	w := new(Wallet)
	w.transactions = make(map[string]*factoid.Transaction)
	db := NewMapDB()
	w.WalletDatabaseOverlay = db
	err := w.InitWallet()
	if err != nil {
		return nil, err
	}
	return w, nil
}

// Close closes a Factom Wallet Database
func (w *Wallet) Close() error {
	return w.dbo.Close()
}

// AddTXDB allows the wallet api to read from a local transaction cashe.
func (w *Wallet) AddTXDB(t *TXDatabaseOverlay) {
	w.txdb = t
}

func (w *Wallet) TXDB() *TXDatabaseOverlay {
	return w.txdb
}

// GenerateECAddress creates and stores a new Entry Credit Address in the
// Wallet. The address can be reproduced in the future using the Wallet Seed.
func (w *Wallet) GenerateECAddress() (*factom.ECAddress, error) {
	// get the next seed from the db
	seed, err := w.GetNextDBSeed()
	if err != nil {
		return nil, err
	}

	// create the new seed
	newseed := sha512.Sum512(seed)
	a, err := factom.MakeECAddress(newseed[:32])
	if err != nil {
		return nil, err
	}

	// save the new seed and the address in the db
	if err := w.InsertNextDBSeed(newseed[:]); err != nil {
		return nil, err
	}

	if err := w.InsertECAddress(a); err != nil {
		return nil, err
	}

	return a, nil
}

// GenerateFCTAddress creates and stores a new Factoid Address in the Wallet.
// The address can be reproduced in the future using the Wallet Seed.
func (w *Wallet) GenerateFCTAddress() (*factom.FactoidAddress, error) {
	// get the next seed from the db
	seed, err := w.GetNextDBSeed()
	if err != nil {
		return nil, err
	}

	// create the new seed
	newseed := sha512.Sum512(seed)
	a, err := factom.MakeFactoidAddress(newseed[:32])
	if err != nil {
		return nil, err
	}

	// save the new seed and the address in the db
	if err := w.InsertNextDBSeed(newseed[:]); err != nil {
		return nil, err
	}

	if err := w.InsertFCTAddress(a); err != nil {
		return nil, err
	}

	return a, nil
}

// GetAllAddresses retrieves all Entry Credit and Factoid Addresses from the
// Wallet Database.
func (w *Wallet) GetAllAddresses() ([]*factom.FactoidAddress, []*factom.ECAddress, error) {
	fcs, err := w.GetAllFCTAddresses()
	if err != nil {
		return nil, nil, err
	}
	ecs, err := w.GetAllECAddresses()
	if err != nil {
		return nil, nil, err
	}

	return fcs, ecs, nil
}

// GetSeed returns the string representaion of the Wallet Seed. The Wallet Seed
// can be used to regenerate the Factoid and Entry Credit Addresses previously
// generated by the wallet. Note that Addresses that are imported into the
// Wallet cannot be regenerated using the Wallet Seed.
func (w *Wallet) GetSeed() (string, error) {
	seed, err := w.GetDBSeed()
	if err != nil {
		return "", err
	}

	return SeedString(seed), nil
}

func (w *Wallet) GetVersion() string {
	return Version
}

func (w *Wallet) GetApiVersion() string {
	return ApiVersion
}
