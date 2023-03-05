terraform {
  required_providers {
    render = {
      version = "1.0.1"
      source  = "render.com/terraform/render"
    }
  }
}

provider "render" {
  # api_key = "<YOUR_API_KEY>"    // Uses env.RENDER_API_KEY, if not supplied
  # email = "<YOUR_RENDER_EMAIL>" // Uses env.RENDER_EMAIL, if not supplied
}

resource "render_service" "mongodb" {
  name = "mongodb"
  repo = "https://github.com/render-examples/mongodb"
  type = "private_service"

  private_service_details = {
    env  = "docker"
    disk = {
      name       = "db"
      mount_path = "/data/db"
      size_gb    = 10
    }
  }
}