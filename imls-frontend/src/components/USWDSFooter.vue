<script>
import { RouterLink } from 'vue-router'

</script>
<template>
  <footer class="usa-footer site-footer">
    <div class="grid-container usa-footer__return-to-top">
      <a href="#">Return to top</a>
    </div>
    <div class="usa-footer__secondary-section bg-base-lightest">
      <div class="grid-container footer-content">
        <h2>About the Public Library Wifi Estimator</h2>
        <p class="maxw-tablet">
          This tool is the result of a collaborative research project with the Institute of Museum and Library Services, the Georgia Public Library Service, and Technology Transformation Services' 10x program.
        </p>
        <div class="grid-row footer-contact-links margin-bottom-4">
          <div class="tablet:grid-col">
            <div class="usa-media-block">
              <div class="usa-media-block__body">
                <h3>Have an idea or an issue?</h3>
                <a target="_blank" href="https://github.com/imls/estimating-wifi/issues/new/choose">
                  Reach out in GitHub
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </footer>
  <div class="usa-identifier site-identifier ">
    <section class="usa-identifier__section usa-identifier__section--masthead" aria-label="Agency identifier">
      <div class="usa-identifier__container">
        <div class="usa-identifier__logos">
            <!-- <a href="javascript:void(0);" class="usa-identifier__logo"> -->
              <svg class="usa-identifier__logo-img usa-icon width-6 height-6" aria-hidden="true" focusable="false" role="img">
                <use xlink:href="~uswds/img/sprite.svg#wifi"></use>
              </svg>
            <!-- </a> -->
          </div>
        <div class="usa-identifier__identity" aria-label="Agency description">
          <p class="usa-identifier__identity-domain">domain.is.tbd</p>
          <p class="usa-identifier__identity-disclaimer">
              This website is NOT an official website of the US Government.
          </p>
        </div>
      </div>
    </section>
    <nav class="usa-identifier__section usa-identifier__section--required-links" aria-label="Important links">
        <div class="usa-identifier__container">
        <ul class="usa-identifier__required-links-list">
          <li class="usa-identifier__required-links-item">
            <a href="https://imls.gov/about/mission" class="usa-identifier__required-link">
                About IMLS
            </a>
          </li>
          <li class="usa-identifier__required-links-item">
            <a href="https://georgialibraries.org/about-gpls/" class="usa-identifier__required-link">
                About GPLS
            </a>
          </li>
          <li class="usa-identifier__required-links-item">
            <a href="https://10x.gsa.gov/process/" class="usa-identifier__required-link">
                About 10x
            </a>
          </li>
          <!-- TODO: update with required links, add a privacy policy? -->
        </ul>
      </div>
    </nav>

  </div>
</template>