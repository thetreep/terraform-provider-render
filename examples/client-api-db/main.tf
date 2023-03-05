terraform {
  required_providers {
    render = {
      version = "1.1.1"
      source  = "render.com/terraform/render"
    }
  }
}

provider "render" {
  # api_key = "<YOUR_API_KEY>"    // Uses env.RENDER_API_KEY, if not supplied
  # email = "<YOUR_RENDER_EMAIL>" // Uses env.RENDER_EMAIL, if not supplied
}

resource "render_service" "client" {
  name = "client"
  repo = "https://github.com/render-examples/nextjs-hello-world"
  type = "static_site"

  static_site_details = {
    build_command = "yarn; yarn build; yarn next export"
    publish_path  = "out"
  }
}

resource "render_service" "api" {
  name = "api"
  repo = "https://github.com/render-examples/hapi-quick-start"
  type = "web_service"

  web_service_details = {
    env = "node"

    native = {
      build_command = "npm install"
      start_command = "node server.js"
    }
  }
}

resource "render_service" "db" {
  name = "db"
  repo = "https://github.com/render-examples/mongodb"
  type = "private_service"

  private_service_details = {
    env = "docker"
    disk = {
      name       = "db"
      mount_path = "/data/db"
      size_gb    = 10
    }
  }
}

resource "render_service_environment" "api" {
  service = render_service.api.id

  variables = [{
    key   = "DATABASE_URL"
    value = render_service.db.private_service_details.url
  }]
}


resource "render_service_environment" "client" {
  service = render_service.client.id

  variables = [
    {
      key   = "API_URL"
      value = render_service.api.web_service_details.url
    }
  ]
}
