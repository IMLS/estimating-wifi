# IMLS Estimating Wifi

[![woke](https://github.com/IMLS/estimating-wifi/actions/workflows/woke.yml/badge.svg)](https://github.com/IMLS/estimating-wifi/actions/workflows/woke.yml)
[![CodeQL](https://github.com/IMLS/estimating-wifi/actions/workflows/codeql.yml/badge.svg)](https://github.com/IMLS/estimating-wifi/actions/workflows/codeql.yml)
[![snyk](https://github.com/IMLS/estimating-wifi/actions/workflows/snyk.yml/badge.svg)](https://github.com/IMLS/estimating-wifi/actions/workflows/snyk.yml)

Estimating Wifi is a pilot project to automate the monitoring of proximal wifi devices via low-cost and open source tools.

This project provides:

- [imls-backend](imls-backend/README.md): PostGREST API over a PostGres database for storing & querying sensor data
- [imls-frontend](imls-frontend/README.md): NodeJS / VueJS web app for displaying the reported data in visualizations
- [imls-wifi-sensor](imls-wifi-sensor/README.md): Command line programs to gather and analyze wifi sessions using Raspberry Pi and Windows machines as sensors
- [imls-windows-installer](imls-windows-installer/README.md): Scripts to generate a Windows installer

Installation instructions for each component are in the links referenced above.

## Privacy

No PII (personally identifiable information) is logged as part of this project. We believe it is impossible to use the data collected to identify an individual.

## About

This software was developed with funding support from [10x](https://10x.gsa.gov/) and in collaboration with the [Institute of Museum and Library Services](https://imls.gov/).

Questions can be directed to 10x at gsa dot gov.

This software is in the **public domain**. No promises are made regarding its functionality or fitness. See the [LICENSE](./LICENSE.md) for more information.

For the full history of this project, please see the [previous repository](https://github.com/18F/imls-pi-stack/).