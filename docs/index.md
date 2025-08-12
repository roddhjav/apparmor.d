---
title: AppArmor.d
hide:
  - toc
---

<!-- Additional styles for landing page -->
<style>
  /* Apply box shadow on smaller screens that don't display tabs */
  @media only screen and (max-width: 1220px) {
    .md-header {
      box-shadow: 0 0 .2rem rgba(0, 0, 0, .1), 0 .2rem .4rem rgba(0, 0, 0, .2);
      transition: color 250ms, background-color 250ms, box-shadow 250ms;
    }
  }

  /* Hide the edit button */
  .md-typeset .md-content__button {
    display: none;
  }

  /* Hide the date of revision */
  .md-source-file {
    display: none;
  }

  /* Get started button */
  .md-typeset .md-button--primary {
    color: var(--md-primary-fg-color);
    background-color: var(--md-primary-bg-color);
    border-color: var(--md-primary-bg-color);
  }

  .md-typeset .md-button--primary:hover {
    color: var(--md-primary-bg-color);
    background-color: var(--md-primary-fg-color);
    border-color: var(--md-primary-bg-color);
  }

  .tx-hero {
    max-width: 700px;
    display: flex;
    padding: .4rem;
    margin: 0 auto;
    text-align: center;
  }

  .tx-hero h1 {
    font-weight: 700;
    font-size: 38px;
    line-height: 46px;
  }

  .tx-hero p {
    color: var(--md-primary-bg-color--light);
    font-weight: 400;
    font-size: 20px;
    line-height: 32px;
  }

  .tx-hero__image {
    max-width: 1350px;
    min-width: 600px;
    width: 100%;
    height: auto;
    margin: 0 auto;
    display: flex;
    align-items: stretch;
  }

  .tx-hero__image img {
    width: 100%;
    height: 100%;
    min-width: 0;
  }

  .image-wrapper img {
    width: 100%;
    height: 100%;
    min-width: 0;
  }

  .main_logo {
    fill: var(--md-primary-bg-color);
    width: 30%;
  }

</style>

<div class="md-container tx-hero">
  <div class="md-grid md-typeset">
    <div class="md-main__inner">
      <div>
        <img class="main_logo" src="assets/avatar-icon.png" alt="" draggable="false">
        <h1>apparmor.d</h1>
        <p><b>Full set of AppArmor policies</b></p>
        <p><code>apparmor.d</code> is a collection of AppArmor profiles designed to restrict the behavior of Linux applications and processes.</p>
        <p>Its goal is to confine everything, targeting both desktops and servers across all distributions that support AppArmor.</p>
        <a href="/overview/"
          title="Get Started" class="md-button md-button--primary">
          Get started
          <svg width="11" height="10" viewBox="0 0 11 10" fill="none" style="margin-left:2px"><path d="M1 5.16772H9.5M9.5 5.16772L6.5 1.66772M9.5 5.16772L6.5 8.66772" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg>
        </a>
        <a href="https://play.pujol.io/" title="Demo Server" class="md-button md-button--primary">
          Demo Server
          <svg height="12" width="12" viewBox="0 0 512 512"><path fill="currentColor" d="M320 0c-17.7 0-32 14.3-32 32s14.3 32 32 32l82.7 0L201.4 265.4c-12.5 12.5-12.5 32.8 0 45.3s32.8 12.5 45.3 0L448 109.3l0 82.7c0 17.7 14.3 32 32 32s32-14.3 32-32l0-160c0-17.7-14.3-32-32-32L320 0zM80 32C35.8 32 0 67.8 0 112L0 432c0 44.2 35.8 80 80 80l320 0c44.2 0 80-35.8 80-80l0-112c0-17.7-14.3-32-32-32s-32 14.3-32 32l0 112c0 8.8-7.2 16-16 16L80 448c-8.8 0-16-7.2-16-16l0-320c0-8.8 7.2-16 16-16l112 0c17.7 0 32-14.3 32-32s-14.3-32-32-32L80 32z" /></svg>
        </a>
      </div>
    </div>
  </div>
</div>
