## Terraform AppStream 2.0 provider
# terraform-provider-appstream


# Provider usage

```
$ go build -o terraform-provider-appstream
$ terraform init
$ terraform plan
$ terraform apply
```

#Development notes
Several other terraform provider projects have been used to refrence how a module should be written,
The goal of this version is to be able to run properly with Terraform Cloud and Terraform Enterprise.
Along side with removing the need for access and secret key in variables and only pass the necessary
to be assumed.

Large portions of code for authentication in config.go & provider.go is from:
https://github.com/terraform-providers/terraform-provider-aws




## Authors/Contributors/Forks
This fork is by:
[Arnvid Karstad](https://github.com/arnvid)

Original code from:
https://github.com/ops-guru/terraform-provider-appstream
[Viktor Berlov](https://github.com/vktr-brlv)

Other forks ref'd:
https://github.com/bluesentry/terraform-provider-appstream
[Chris Mackubin](https://github.com/chris-mackubin)
