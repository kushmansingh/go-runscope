package runscope

import (
	"testing"
	"encoding/json"
	"time"
)


func TestCreateEnvironment(t *testing.T) {
	testPreCheck(t)
	client := clientConfigure()
	bucket, err := client.CreateBucket(Bucket{Name: "test", Team: Team{Id: teamId}})
	defer client.DeleteBucket(bucket.Key)

	if err != nil {
		t.Error(err)
	}

	environment := &Environment{
		Name: "tf_environment",
		InitialVariables: map[string]string{
			"VarA" : "ValB",
			"VarB" : "ValB",
		},
		Integrations: []Integration{
			{
				Id: "27e48b0d-ba8e-4fe0-bcaa-dd9de08dc47d",
				IntegrationType: "pagerduty",
			},
			{
				Id: "574f4560-0f50-41da-a2f7-bdce419ad378",
				IntegrationType: "slack",
			},
		 },
	}

	environment, err = client.CreateSharedEnvironment(environment, bucket)
	defer client.DeleteSharedEnvironment(environment, bucket)

	if err != nil {
		t.Error(err)
	}

	if len(environment.Id) == 0 {
		t.Error("Environment id should not be empty")
	}

	if len(environment.InitialVariables) != 2 {
		t.Errorf("Expected %d initial variables got %d", 2, len(environment.InitialVariables))
	}
}


func TestReadEnvironmentFromResponse(t *testing.T) {
	responseMap := new(response)
	if err := json.Unmarshal([]byte(sampleEnvironment), &responseMap); err != nil {
		t.Error(err)
	}

	environment, err := getEnvironmentFromResponse(responseMap.Data)
	if err != nil {
		t.Error(err)
	}

	if environment.Name != "Production" {
		t.Errorf("Expected name %s, actual %s", "Production", environment.Name)
	}

	if environment.Script != "var s = 'foo'" {
		t.Errorf("Expected script %s, actual %s", "var s = 'foo'", environment.Script)
	}

	if !environment.PreserveCookies {
		t.Errorf("Expected PreserveCookies %s, actual %s", true, environment.PreserveCookies)
	}

	if environment.TestId != "a10c97e6-2024-41ca-990d-5e0b5f751734" {
		t.Errorf("Expected test id %s, actual %s", "a10c97e6-2024-41ca-990d-5e0b5f751734", environment.Script)
	}

	if len (environment.InitialVariables) != 2 {
		t.Errorf("Expected %s initial variables, actual %s", "2", len(environment.InitialVariables))
	}

	if environment.InitialVariables["NameB"] != "ValueB" {
		t.Errorf("Expected initial variable value %s, actual %s", "ValueB", environment.InitialVariables["NameB"] )
	}

	if len (environment.Integrations) != 2 {
		t.Errorf("Expected %s integrations, actual %s", "2", len(environment.Integrations))
	}

	if environment.Integrations[1].Id != "1b766ead-b3d1-456f-a350-83845a428ed1" ||
	   environment.Integrations[1].Description != "PagerDuty: Runscope Service" ||
	   environment.Integrations[1].IntegrationType != "pagerduty" {
		t.Errorf("Expected integration not correct got #%v", environment.Integrations[1])
	}

	if environment.Id != "c392d38e-70df-4181-abe5-51864ccf8f23" {
		t.Errorf("Expected id %s, actual %s", "c392d38e-70df-4181-abe5-51864ccf8f23", environment.Id)
	}

	if len (environment.Regions) != 2 {
		t.Errorf("Expected %s regions, actual %s", "2", len(environment.Regions))
	}

	if !environment.VerifySsl {
		t.Errorf("Expected verify ssl %s, actual %s", true, environment.VerifySsl)
	}

	expectedTime := time.Time{}
	expectedTime = time.Unix(int64(1494190571), 0)
	if !environment.ExportedAt.Equal(expectedTime) {
		t.Errorf("Expected exported at %s, actual %s",  expectedTime.String(), environment.ExportedAt)
	}

	if !environment.RetryOnFailure {
		t.Errorf("Expected retry on failures %s, actual %s", true, environment.RetryOnFailure)
	}

	if len (environment.RemoteAgents) != 1 ||
		environment.RemoteAgents[0].Name != "my-local-machine.runscope.com" ||
		environment.RemoteAgents[0].Uuid != "141d4dbc-1e41-401e-8067-6df18501e9ed" {
		t.Errorf("Expected remote agent not correct got #%v", environment.RemoteAgents[0])
	}

	if len (environment.WebHooks) != 2 ||
		environment.WebHooks[1] != "https://yourapihere.com/post" {
		t.Errorf("Expected web hooks are not correct got #%v", environment.WebHooks)
	}

	if environment.ParentEnvironmentId != "8ace1fbb-9626-4455-b006-116ba7154c1c" {
		t.Errorf("Expected parent environment id %s, actual %s", "8ace1fbb-9626-4455-b006-116ba7154c1c", environment.ParentEnvironmentId)
	}

	if !environment.EmailSettings.NotifyAll ||
		environment.EmailSettings.NotifyOn != "all" ||
		environment.EmailSettings.NotifyThreshold != 4 ||
		len (environment.EmailSettings.Recipients) != 1 ||
		environment.EmailSettings.Recipients[0].Id != "4ee15ecc-7fe1-43cb-aa12-ef50420f2cf9" {
		t.Errorf("Expected email settings not correct got #%v", environment.EmailSettings)
	}
}

const sampleEnvironment string = `
{
  "meta": {
    "status": "success"
  },
  "data": {
    "script_library": [],
    "name": "Production",
    "script": "var s = 'foo'",
    "preserve_cookies": true,
    "test_id": "a10c97e6-2024-41ca-990d-5e0b5f751734",
    "initial_variables": {
      "NameA": "ValueA",
      "NameB": "ValueB"
    },
    "integrations": [
      {
        "integration_type": "slack",
        "description": "Slack: Technology channel, send message on failed test runs",
        "id": "a9fa014e-5dc0-4d87-8638-3f696a381062"
      },
      {
        "integration_type": "pagerduty",
        "description": "PagerDuty: Runscope Service",
        "id": "1b766ead-b3d1-456f-a350-83845a428ed1"
      }
    ],
    "auth": null,
    "id": "c392d38e-70df-4181-abe5-51864ccf8f23",
    "regions": [
      "us1",
      "jp1"
    ],
    "headers": {},
    "verify_ssl": true,
    "version": "1.0",
    "exported_at": 1494190571,
    "retry_on_failure": true,
    "remote_agents": [
        {
            "name": "my-local-machine.runscope.com",
            "uuid": "141d4dbc-1e41-401e-8067-6df18501e9ed"
        }
    ],
    "webhooks": [
        "http://api.example.com/webhook_reciever",
        "https://yourapihere.com/post"
    ],
    "parent_environment_id": "8ace1fbb-9626-4455-b006-116ba7154c1c",
    "stop_on_failure": true,
    "emails": {
        "notify_all": true,
        "notify_on": "all",
        "notify_threshold": 4,
        "recipients": [
            {
                "email": "grace@example.com",
                "name": "Grace Hopper",
                "id": "4ee15ecc-7fe1-43cb-aa12-ef50420f2cf9"
            }
        ]
    },
    "client_certificate": ""
  }
}
`