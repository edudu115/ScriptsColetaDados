package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
)

func quickSight(config aws.Config, accountID string) ([]byte, error) {
	clientQuickSight := quicksight.NewFromConfig(config)

	datasetsOutput, err := clientQuickSight.ListDataSets(context.TODO(), &quicksight.ListDataSetsInput{
		AwsAccountId: &accountID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar os datasets: %w", err)
	}

	data := make([]map[string]interface{}, 0)

	for _, dataset := range datasetsOutput.DataSetSummaries {
		ingestionsOutput, err := clientQuickSight.ListIngestions(context.TODO(), &quicksight.ListIngestionsInput{
			AwsAccountId: &accountID,
			DataSetId:    dataset.DataSetId,
		})
		if err != nil {
			return nil, fmt.Errorf("erro ao listar as ingestions: %w", err)
		}

		for _, ingestion := range ingestionsOutput.Ingestions {
			data = append(data, map[string]interface{}{
				"nameDataset": *dataset.Name,
				"idIngestion": *ingestion.IngestionId,
				"status":      ingestion.IngestionStatus,
				"createdTime": ingestion.CreatedTime.String(),
			})
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao converter para JSON: %w", err)
	}

	return jsonData, nil
}

func main() {
	if len(os.Args) < 6 {
		log.Fatalf("Uso: %s <accessKey> <secretKey> <region> <service> <accountID>", os.Args[0])
	}

	accessKey := os.Args[1]
	secretKey := os.Args[2]
	regionAws := os.Args[3]
	service := os.Args[4]
	accountID := os.Args[5]

	config := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      regionAws,
	}

	if service == "quicksight" {
		jsonData, err := quickSight(config, accountID)
		if err != nil {
			log.Fatalf("Erro: %v", err)
		}
		fmt.Println("{ \"data\": " + string(jsonData) + " }")
	} else {
		log.Fatalf("Serviço não suportado: %s", service)
	}
}