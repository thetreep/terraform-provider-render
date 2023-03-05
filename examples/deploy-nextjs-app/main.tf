terraform {
  required_providers {
    render = {
      version = "0.1.0"
      source  = "render.com/terraform/render"
    }
  }
}

provider "render" {
  # api_key = "<YOUR_API_KEY>"    // Uses env.RENDER_API_KEY, if not supplied
  # email = "<YOUR_RENDER_EMAIL>" // Uses env.RENDER_EMAIL, if not supplied
}

resource "render_service" "nextjs" {
  name = "nextjs-static"
  repo = "https://github.com/render-examples/nextjs-hello-world"
  type = "static_site"

  static_site_details = {
    build_command = "yarn; yarn build; yarn next export"
    publish_path  = "out"
  }
}