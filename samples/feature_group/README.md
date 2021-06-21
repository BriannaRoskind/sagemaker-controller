# Feature Group Sample

This sample demonstrates how to create a feature group using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.

## Prerequisites

This sample assumes that you have completed the [common prerequisites](https://github.com/aws-controllers-k8s/sagemaker-controller/blob/main/samples/README.md).

## Create an S3 bucket:

Since we are using the offline store in this example, you need to set up an s3 bucket. [Here are directions](https://docs.aws.amazon.com/AmazonS3/latest/userguide/create-bucket-overview.html) to set up your s3 bucket through the S3 Console, AWS SDK, or AWS CLI.

## Adding required policies to your IAM role:

First create an Amazon SageMaker execution role, and then [here are directions](https://docs.aws.amazon.com/sagemaker/latest/dg/feature-store-adding-policies.html) to give your SageMaker Execution Role SageMakerFeatureStoreAccess. Make sure that your role has access to your s3 bucket, this can be done by giving it AmazonS3FullAccess.

## Creating your Feature Group

### Create a Feature Group:

To submit your prepared feature group specification, apply the specification to your Kubernetes cluster as such:


```
`$ kubectl apply ``-``f ``my``-``feature-group``.``yaml``
featuregroup.sagemaker.services.k8s.aws/my-feature-group created`
```

### List Feature Groups:

To list all feature groups created using the ACK controller use the following command:


```
`$ kubectl get featuregroup`
```

### Describe a Feature Group:

To get more details about the feature group once it's submitted, like checking the status, errors or parameters of the feature group use the following command:


```
`$ kubectl describe featuregroup my-feature-group created`
```

## Ingesting Data into your Feature Group

Note: Assumes Creation of a feature group with its name stored in `feature_group_name`

```
### Sample CSV data file for Ingestion Example:
#TransactionID,EventTime
#1,1623434915
#2,1623435267
#3,1623435284

###Example boto3 ingestion of feature group:

import boto3
import csv

sagemaker_featurestore_runtime_client = boto3.Session().client(
    service_name="sagemaker-featurestore-runtime")

### OPTION 1: To Download all records at once and upload records sequentially
with open('./Downloads/Sample_data.csv') as file_handle:
    records =[
                   [
		                    {'FeatureName':featureName,
				                      'ValueAsString':valueAsString}
						                     for featureName, valueAsString in row.items()]
								                  for row in csv.DictReader(file_handle, skipinitialspace=True)]

for record in records:
    sagemaker_featurestore_runtime_client.put_record(
            FeatureGroupName=feature_group_name,
	            Record=record)

### OPTION 2: To Download records sequentially and upload records sequentially
with open('./Downloads/Sample_data.csv') as file_handle:
    for row in csv.DictReader(file_handle, skipinitialspace=True):
            record =[{'FeatureName':featureName,
	                'ValueAsString':valueAsString
			            }
				                for featureName, valueAsString in row.items()
						        ]
							        sagemaker_featurestore_runtime_client.put_record(
								            FeatureGroupName=feature_group_name,
									                Record=record)

# To Check that the records are retrievable
for recordIdentifierValue in range(1,len(records) + 1):
    sagemaker_featurestore_runtime_client.get_record(
           FeatureGroupName=feature_group_name,
	          RecordIdentifierValueAsString=str(recordIdentifierValue))
		  ```

## Deleting your Feature Group

To delete the feature group, use the following command:

```
`$ kubectl ``delete`` featuregroup ``my``-feature-group``
featuregroup.sagemaker.services.k8s.aws "my-feature-group" deleted
`
```

