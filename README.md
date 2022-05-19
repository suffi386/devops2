<img src="./docs/static/logos/zitadel-logo-dark@2x.png#gh-light-mode-only" alt="Zitadel Logo" height="100px" width="auto" />
<img src="./docs/static/logos/zitadel-logo-light@2x.png#gh-dark-mode-only" alt="Zitadel Logo" height="100px" width="auto" />

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Release](https://github.com/zitadel/zitadel/actions/workflows/zitadel.yml/badge.svg)](https://github.com/zitadel/zitadel/actions)
[![license](https://badgen.net/github/license/zitadel/zitadel/)](https://github.com/zitadel/zitadel/blob/main/LICENSE)
[![release](https://badgen.net/github/release/zitadel/zitadel/stable)](https://github.com/zitadel/zitadel/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/zitadel/zitadel)](https://goreportcard.com/report/github.com/zitadel/zitadel)
[![codecov](https://codecov.io/gh/zitadel/zitadel/branch/main/graph/badge.svg)](https://codecov.io/gh/zitadel/zitadel)
[![discord](https://badgen.net/discord/online-members/erh5Brh7jE)](https://discord.gg/erh5Brh7jE)  
[![OpenID Connect Certified](./docs/static/logos/oidc-cert.png)](https://openid.net/certification/#OPs)

---

ZITADEL is the modern open* alternative for Auth0, Firebase Auth, AWS Cognito as well as Keycloak built for the container and serverless era.

We provide a wide range of out of the box features like secure login, self-service, OpenID Connect, OAuth2.x, SAML2, branding, Passwordless with FIDO2, OTP, U2F, and an unlimited audit trail to improve the life of developers. Especially noteworthy is that ZITADEL supports not only B2C and B2E scenarios but also B2B. This is super useful for people who build B2B Solutions, as ZITADEL can handle all the delegated user and access management.

With ZITADEL you rely on a battle tested, hardened and extensible turnkey solution to solve all of your authentication and authorization needs. With the unique way of how ZITADEL stores data it gives you an unlimited audit trail which provides a peace of mind for even the harshest audit and analytics requirements.

<!-- TODO: Insert Video here-->

## Get started

### ZITADEL Cloud

The easiest way to get started with ZITADEL is to use our public cloud offering. [Subscribe to our newsletter](https://zitadel.com/v2) and we will be in touch with you as soon as the public release is live.

You can also discovery our new pay-as-you-go [pricing](https://zitadel.com/pricing/v2).

### Install ZITADEL

- [We provide installation guides for multiple platforms here](https://docs.zitadel.com/docs/guides/installation)

### Quickstarts - Integrate you app

- [Multiple Quickstarts can be found here](https://docs.zitadel.com/docs/quickstarts/introduction)
- [And even more examples are located under](https://github.com/zitadel/zitadel-examples)

> If you miss something please feel free to engage with us [here](https://github.com/zitadel/zitadel/discussions/1717)

## Why ZITADEL

- [API-first](https://docs.zitadel.com/docs/apis/introduction)
- Strong audit trail thanks to [event sourcing](https://docs.zitadel.com/docs/concepts/eventstore)
- [Actions](https://docs.zitadel.ch/docs/concepts/features/actions) to react on events with custom code
- [Private labeling](https://docs.zitadel.com/docs/guides/customization/branding) for a uniform user experience
- [cockroach database](https://www.cockroachlabs.com/) is the only dependency

## Features

- Single Sign On (SSO)
- Passwordless with FIDO2 support
- Username / Password
- Multifactor authentication with OTP, U2F
- [Identity Brokering](https://docs.zitadel.com/docs/guides/authentication/identity-brokering)
- [Machine-to-machine (JWT profile)](https://docs.zitadel.com/docs/guides/authentication/serviceusers)
- [Personal Access Tokens (PAT)]()
- Role Based Access Control (RBAC)
- Delegate role management to third-parties
- Self-registration including verification
- User self service
- [Service Accounts](https://docs.zitadel.com/docs/guides/authentication/serviceusers)

### Client libraries

<!-- TODO: check other libraries -->

| Language | Client | API | Machine auth (\*) | Auth check (\*\*) | Thanks to the maintainers |
|----------|--------|--------------|----------|---------|---------------------------|
| .NET     | [zitadel-net](https://github.com/zitadel/zitadel-net) | GRPC | ✔️ | ✔️ | [buehler 👑](https://github.com/buehler) |
| Dart     | [zitadel-dart](https://github.com/zitadel/zitadel-dart) | GRPC | ✔️ | ❌ | [buehler 👑](https://github.com/buehler) |
| Elixir   | [zitadel_api](https://github.com/jshmrtn/zitadel_api) | GRPC | ✔️ | ✔️ | [jshmrtn 🙏🏻](https://github.com/jshmrtn) |
| Go       | [zitadel-go](https://github.com/zitadel/zitadel-go) | GRPC | ✔️ | ✔️ | ZITADEL |
| Rust     | [zitadel-rust](https://crates.io/crates/zitadel) | GRPC | ✔️ | ❌ | [buehler 👑](https://github.com/buehler) |
| JVM      | ❓ | ❓ | ❓ | | Maybe you? |
| Python   | ❓ | ❓ | ❓ | | Maybe you? |
| Javascript | ❓ | ❓ | ❓ | | Maybe you? |

(\*) Automatically authenticate service accounts with [JWT Profile](https://docs.zitadel.com/docs/apis/openidoauth/grant-types#json-web-token-jwt-profile).  
(\*\*) Automatically check if the access token is valid and claims match

## How To Contribute

Details about how to contribute you can find in the [Contribution Guide](./CONTRIBUTING.md)

## Security

See the policy [here](./SECURITY.md)

## Showcase

<!-- TODO: Replace Images-->

### Passwordless Login

Use our login widget to allow easy and sucure access to your applications and enjoy all the benefits of passwordless (FIDO 2 / WebAuthN):

* works on all modern platforms, devices, and browsers
* phishing resistant alternative
* requires only one gesture by the user
* easy [enrollment](https://docs.zitadel.com/docs/manuals/user-factors) of the device during registration

![passwordless-windows-hello](https://user-images.githubusercontent.com/1366906/118765435-5d419780-b87b-11eb-95bf-55140119c0d8.gif)

### Admin Console

Use [Console](https://docs.zitadel.com/docs/manuals/introduction) or our [APIs](https://docs.zitadel.com/docs/apis/introduction) to setup organizations, projects and applications.

Register new applications
![OIDC-Client-Register](https://user-images.githubusercontent.com/1366906/118765446-62064b80-b87b-11eb-8b24-4f4c365b8c58.gif)

Delegate the right to assign roles to another organization
![projects_create_org_grant](https://user-images.githubusercontent.com/1366906/118766069-39cb1c80-b87c-11eb-84cf-f5becce4e9b6.gif)

Customize login and console with your design  
![private_labeling](https://user-images.githubusercontent.com/1366906/123089110-d148ff80-d426-11eb-9598-32b506f6d4fd.gif)

## Usage Data

ZITADEL components send errors and usage data to CAOS Ltd., so that we are able to identify code improvement potential. If you don't want to send this data or don't have an internet connection, pass the global flag `--disable-analytics` when using zitadelctl. For disabling ingestion for already-running components, execute the takeoff command again with the `` flag.

We try to distinguishing the environments from which events come from. As environment identifier, we enrich the events by the domain you have configured in zitadel.yml, as soon as it's available. When it's not available and you passed the --gitops flag, we defer the environment identifier from your git repository URL.

Besides from errors that don't clearly come from misconfiguration or cli misuage, we send an inital event when any binary is started. This is a "<component> invoked" event along with the flags that are passed to it, except secret values of course.

We only ingest operational data. Your ZITADEL workload data from the IAM application itself is never sent anywhere unless you chose to integrate other systems yourself.

## License

See the exact licensing terms [here](./LICENSE)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
