stack "just_consul" {
  consul {
    address = "http://localhost:8500"
    namespace = "default"
  }
}