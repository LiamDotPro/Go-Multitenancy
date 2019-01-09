# Go-Multitenancy
A Golang Multi Tenancy Framework

This is the backend framework for multi tenancy, this project is currently under development.

The goal of this framework is to provide a simple and fast soloution for starting a SaaS based web or mobile soloution that is run in go.

Supported Databases:
- [X] Postgresql
- [X] MySql
- [X] Tsql
- [ ] Provide enviroment variables based soloution for switch between adapters.
- [ ] NoSql support (Long Term Goal)

Currently most SQL databases are supported through the change of the GORM config.

Backend Todo:
- [ ] Intergrate Stripe Payment soloution
- [ ] User Module Creation & Testing
- [ ] Package down soloution and refactor code to work with module's to minimize required code.
- [ ] Add CI
- [ ] Create Cli Project for Creating New Project
