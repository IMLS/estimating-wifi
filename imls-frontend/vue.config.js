import process from 'process';

module.exports = {
    publicPath: process.env.NODE_ENV === 'production'
      // process.env.BASEURL should be '/site/[ORG_NAME]/[REPO_NAME]' on federalist
      ? process.env.BASEURL
      : '/'
  }