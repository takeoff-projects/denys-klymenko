terraform {
  backend "gcs" {
    prefix = "terraform-state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.provider_region
}

provider "google-beta" {
  project = var.project_id
  region  = var.provider_region
}

variable "gcp_service_list" {
  description = "The list of apis necessary for the project"
  type        = list(string)
  default = [
    "cloudresourcemanager.googleapis.com",
    "serviceusage.googleapis.com",
    "iam.googleapis.com",
    "cloudfunctions.googleapis.com",
    "cloudbuild.googleapis.com",
    "servicemanagement.googleapis.com",
    "servicecontrol.googleapis.com",
    "appengine.googleapis.com",
    "firestore.googleapis.com",
    "secretmanager.googleapis.com",
    "run.googleapis.com",
    "endpoints.googleapis.com",
    "apigateway.googleapis.com"
  ]
}

resource "google_project_service" "gcp_services" {
  for_each           = toset(var.gcp_service_list)
  project            = var.project_id
  service            = each.key
  disable_on_destroy = false
}

resource "google_app_engine_application" "app" {
  project       = var.project_id
  location_id   = "us-central"
  database_type = "CLOUD_FIRESTORE"
}

resource "google_cloud_run_service" "default" {
  name     = "pets-api"
  location = var.provider_region
  project = var.project_id

  autogenerate_revision_name = true

  template {
    spec {
      containers {
        image = "gcr.io/roi-takeoff-user5/pets-api:2"
        ports {
          container_port = 8080
        }

        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = var.project_id
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }

  depends_on = [google_cloud_run_service.default]
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.default.location
  project  = var.project_id
  service  = google_cloud_run_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
  depends_on = [google_cloud_run_service.default]
}

resource "google_api_gateway_api" "default" {
  provider = google-beta
  api_id   = "pets-api"
  project  = var.project_id
}

resource "google_api_gateway_api_config" "default" {
  project       = var.project_id
  provider      = google-beta
  api           = google_api_gateway_api.default.api_id

  openapi_documents {
    document {
      path = "openapi.yaml"
      contents = base64encode(templatefile(
        "${path.root}/../docs/swagger.yaml",
        {
          backend : google_cloud_run_service.default.status[0].url
        }
      ))
    }
  }

  depends_on = [google_cloud_run_service.default]

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "default" {
  project    = var.project_id
  region     = var.provider_region
  provider   = google-beta
  api_config = google_api_gateway_api_config.default.id
  gateway_id = "pets-api-gw"

  depends_on = [google_api_gateway_api_config.default]
}


output "public_url" {
  value = google_cloud_run_service.default.status[0].url
}

output "gateway_url" {
  value = google_api_gateway_gateway.default.default_hostname
}