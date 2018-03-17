# Changes

## 0.3.0 (March 17, 2018)

![GitHub downloads](https://img.shields.io/github/downloads/ringods/terraform-provider-xenserver/v0.3.0/total.svg)

FEATURES:

* **New Data Source:** `xenserver_pif` [#29](https://github.com/ringods/terraform-provider-xenserver/pull/29) by @joncave
* **New Data Source:** `xenserver_sr` [#30](https://github.com/ringods/terraform-provider-xenserver/pull/30) by @joncave

ENHANCEMENTS:

* Bootstrapped a static Hugo site (several commits) by @ringods
* resource/xenserver_network: Add support for `other_config` to the network resource [#31](https://github.com/ringods/terraform-provider-xenserver/pull/31) by @joncave
* resource/xenserver_vm: Add support for `other_config` to the VM resource [#35](https://github.com/ringods/terraform-provider-xenserver/pull/35) by @joncave

BUG FIXES:

* resource/xenserver_vm: VDI without VBD causes silent failure to create VM [#21](https://github.com/ringods/terraform-provider-xenserver/issues/21) via [#23](https://github.com/ringods/terraform-provider-xenserver/pull/23) by @briantopping

## 0.2.0 (December 2, 2017)

![GitHub downloads](https://img.shields.io/github/downloads/ringods/terraform-provider-xenserver/v0.2.0/total.svg)

NOTES:

* Initial release after maintenance transferred from `@amfranz` to `@ringods`

FEATURES:

* **New Data Source:** `xenserver_pifs` by @ringods

IMPROVEMENTS:

* Set up Github releases via Travis CI.
* Support for XenServer 7.2 via [`go-xen-api-client` #180cc7b](https://github.com/ringods/go-xen-api-client/commit/180cc7bfb7590fbc1a81c198b0011429ac58881f)
