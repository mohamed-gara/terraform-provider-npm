data "npm_package" "my_pkg" {
  name    = "my-package-npm"
  version = "1.1.0"
}

output "npm_package_files" {
  value = data.npm_package.my_pkg.files
}
