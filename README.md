# url-shortener

Shorten URL and generage QR code.

![url-shortener](https://user-images.githubusercontent.com/5100568/38465564-40d7b79e-3b58-11e8-978f-304193d4f5c3.gif)

## How to start

Edit app.yml for your configuration.

~~~
env_variables:
  BASE_URL: "https://s.juntaki.com"
  RECAPTCHA_SITE_KEY: "6Lc48lAUAAAAAFehXaHNp0Ys-lS1iUfNtUXd_-eR"
~~~

and create secret.yml

~~~
env_variables:
    RECAPTCHA_SECRET: "<your-recaptcha-secret>"
~~~

Run commands to deploy app to GAE.

~~~
gcloud config set project <project-name>
gcloud app deploy
~~~
