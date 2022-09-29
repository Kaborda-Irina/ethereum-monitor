package vault

import (
	"context"
	vault "github.com/hashicorp/vault/api"
	"log"
)

type DataVault struct {
	Ctx        context.Context
	Client     *vault.Client
	MountPath  string
	SecretPath string
}

func InitVault() *DataVault {
	config := vault.DefaultConfig()

	config.Address = "http://127.0.0.1:8200"

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	client.SetToken("dev-only-token")

	dataVault := &DataVault{
		Ctx:        context.Background(),
		Client:     client,
		MountPath:  "secret",
		SecretPath: "privateKey",
	}
	return dataVault

}

func ReadKey(dataVault *DataVault, accountIndex string) string {
	//versions, err := dataVault.Client.KVv2("secret").GetVersionsAsList(dataVault.Ctx, dataVault.SecretPath)
	//if err != nil {
	//	log.Fatalf("unable to read versions list: %v", err)
	//}
	//
	//for _, version := range versions {
	//	deleted := "Not deleted"
	//	if !version.DeletionTime.IsZero() {
	//		deleted = version.DeletionTime.Format(time.UnixDate)
	//	}
	//
	//	secret, err := dataVault.Client.KVv2(accountIndex).GetVersion(dataVault.Ctx, dataVault.SecretPath, version.Version)
	//	if err != nil {
	//		log.Fatalf("unable to read secret: %v", err)
	//	}
	//	value, ok := secret.Data[accountIndex].(string)
	//	if ok {
	//		log.Printf(
	//			"Version: %d. Created at: %s. Deleted at: %s. Destroyed: %t. Value: '%s'.\n",
	//			version.Version,
	//			version.CreatedTime.Format(time.UnixDate),
	//			deleted,
	//			version.Destroyed,
	//			value,
	//		)
	//	}
	//	return value
	//}
	secret, err := dataVault.Client.KVv2(dataVault.MountPath).Get(dataVault.Ctx, dataVault.SecretPath)
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	value, ok := secret.Data[accountIndex].(string)
	if !ok {
		log.Fatalf("value type assertion failed: %T %#v %s %v", secret.Data[accountIndex], secret.Data[accountIndex], accountIndex, secret.Data)
	}

	return value
}

func WriteKey(dataVault *DataVault, secretData map[string]interface{}) error {
	_, err := dataVault.Client.KVv2(dataVault.MountPath).Put(dataVault.Ctx, dataVault.SecretPath, secretData)
	if err != nil {
		log.Fatalf("unable to write secret: %v", err)
	}

	return err
}
