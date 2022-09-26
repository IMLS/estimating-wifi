# IMLS Estimating Wifi

[![MegaLinter](https://github.com/IMLS/estimating-wifi/actions/workflows/megalinter.yml/badge.svg)](https://github.com/IMLS/estimating-wifi/actions/workflows/megalinter.yml)
[![woke](https://github.com/IMLS/estimating-wifi/actions/workflows/woke.yml/badge.svg)](https://github.com/IMLS/estimating-wifi/actions/workflows/woke.yml)
[![CodeQL](https://github.com/IMLS/estimating-wifi/actions/workflows/codeql.yml/badge.svg)](https://github.com/IMLS/estimating-wifi/actions/workflows/codeql.yml)

[![snyk](https://github.com/IMLS/estimating-wifi/actions/workflows/snyk.yml/badge.svg)](https://github.com/IMLS/estimating-wifi/actions/workflows/snyk.yml)

Estimating Wifi is a pilot project to automate the monitoring of wifi usage via low-cost and open source tools.

For the full history of this project, please see the [previous repository](https://github.com/18F/imls-pi-stack/).

This project provides:

- imls-backend: PostGREST API over a PostGres database for storing & querying sensor data
- imls-frontend: NodeJS / VueJS web app for displaying the reported data in visualizations
- imls-wifi-sensor: Command line programs to gather and analyze wifi sessions using Raspberry Pi and Windows machines as sensors
- imls-windows-installer: Scripts to generate a Windows installer
- doc: Developer documentation

## Privacy

No PII (personally identifiable information) is logged as part of this project. We believe it is impossible to use the data collected to identify an individual.

## Installation instructions

Forthcoming.

## About

This software was developed with funding support from [10x](https://10x.gsa.gov/) and in collaboration with the [Institute of Museum and Library Services](https://imls.gov/).

Questions can be directed to 10x at gsa dot gov.

This software is in the **public domain**. No promises are made regarding its functionality or fitness. See the [LICENSE](./LICENSE.md) for more information.
