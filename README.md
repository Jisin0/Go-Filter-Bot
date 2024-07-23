# GoFilterBot [![telegram badge](https://img.shields.io/badge/Support-30302f?style=flat&logo=telegram)](https://telegram.dog/Jisin0) [![Go Report Card](https://goreportcard.com/badge/github.com/Jisin0/Go-Filter-Bot)](https://goreportcard.com/report/github.com/Jisin0/Go-Filter-Bot) [![Go Build](https://github.com/Jisin0/Go-Filter-Bot/workflows/Go/badge.svg)](https://github.com/Jisin0/Go-Filter-Bot/actions?query=workflow%3AGo+event%3Apush+branch%3Amain) [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


A **Simple**, **Superfast** and **Serverless** Filter Bot Written in Go with Global Filters and Autodelete support.
The repository can be hosted on vercel as a serverless application with close to no downsides.

## Commands

```
/start    - Shows a start message                    [Group/DM]
/about    - Shows an about message                   [Group/DM]
/help     - Shows a help message                     [Group/DM]
/filter   - Create a new filter for a word or phrase [Group/Connected DM]
/stop     - Stop an existing filter                  [Group/Connected DM]
/gstop    - Stop an existing global filter           [Admin DM]
/filters  - View all the filters saved for a chat    [Group/Connected DM]
/gfilter  - Create a new global filter that will work in all chats [Admin DM]
/gfilters - View all the saved global filters        [All]
/filter   - Create a new filter for a word or phrase [Group/Connected DM]
/broadcast - Broadcast a message to all bot users     [Admin DM]
/concast  - Broadcast a message to connected users   [Admin DM]
 ```

## Variables

Variables can also be loaded by creating a ```.env``` file at the root of the repository. See [.env.sample](/.env.sample) to see the format to use.

- [X] `BOT_TOKEN` : Bot token obtained by creating a bot from [@BotFather](https://telegram.dog/BotFather).
- [X] `MONGODB_URI` : [MongoDB](https://www.mongodb.com) URI. Get this value from [mongoDB](https://www.mongodb.com). For more help watch this [video](https://youtu.be/1G1XwEOnxxo).
- [X] `ADMINS` : Telegram user ids' of bot admins seperated by spaces.
- [ ] `MULTI_FILTER` : Set to True if multiple filters should be processed for a single message (don't add on serverles).
- [ ] `AUTO_DELETE`: Time in minutes after which a filter result should be automatically deleted. for ex: 60 for 1hour (won't work serverless).
- [ ] `PORT` : The port on which the webapp should run (use 10000 on render)

## Deploy
<details><summary>Deploy To Vercel</summary>
<p>
Follow these instructions to deploy this repo to <b>vercel</b>
<ol type="1">
<li><b>Fork</b> this repository 🍴</li>
<li>Go to your <a href="https://vercel.com">vercel</a> dashboard and create a <b>Add New > Project</b></li>
<li>Fill in the <b>BOT_TOKEN</b> and <b>MONGODB_URI</b> environment variables</li>
<li>Click <b>Deploy</b> and wait</li>
<li>Open your app and put in your bot token and click <b>Connect</b></li>
</ol>
</p>
</details>

<details><summary>Deploy To Heroku</summary>
<p>
<br>
<a href="https://heroku.com/deploy?template=https://github.com/Jisin0/Go-Filter-Bot/tree/main">
  <img src="https://www.herokucdn.com/deploy/button.svg" alt="Deploy">
</a>
</p>
</details>

<details><summary>Deploy To Scalingo</summary>
<p>
<br>
<a href="https://dashboard.scalingo.com/create/app?source=https://github.com/Jisin0/Go-Filter-Bot#main">
   <img src="https://cdn.scalingo.com/deploy/button.svg" alt="Deploy on Scalingo" data-canonical-src="https://cdn.scalingo.com/deploy/button.svg" style="max-width:100%;">
</a>
</p>
</details>

<details><summary>Deploy To Render</summary>
<p>
<br>
<a href="https://dashboard.render.com/select-repo?type=web">
  <img src="https://render.com/images/deploy-to-render-button.svg" alt="deploy-to-render">
</a>
</p>
<p>
Make sure to have the following options set :

<b>Environment</b>
<pre>Go</pre>

<b>Build Command</b>
<pre>go build .</pre>

<b>Start Command</b>
<pre>./Go-Filter-Bot</pre>

<b>Advanced >> Health Check Path</b>
<pre>/</pre>
</p>
</details>


<details><summary>Deploy To Koyeb</summary>
<p>
<br>
<a href="https://app.koyeb.com/deploy?type=git&repository=github.com/Jisin0/Go-Filter-Bot&branch=main">
  <img src="https://www.koyeb.com/static/images/deploy/button.svg" alt="deploy-to-koyeb">
</a>
</p>
<p>
You must set the Run command to :
<pre>./bin/Go-Filter-Bot</pre>
</p>
</details>

<details><summary>Deploy To Okteto</summary>
<p>
<br>
<a href="https://cloud.okteto.com/deploy?repository=https://github.com/Jisin0/Go-Filter-Bot">
  <img src="https://okteto.com/develop-okteto.svg" alt="deploy-to-okteto">
</a>
</p>
</details>

<details><summary>Deploy To Railway</summary>
<p>
<br>
<a href="https://railway.app/new/template?template=https%3A%2F%2Fgithub.com%2FJisin0%2FGo-Filter-Bot">
  <img src="https://railway.app/button.svg" alt="deploy-to-railway">
</a>
</p>
</details>

<details><summary>Run Locally/VPS</summary>
<p>
You must have the latest version of <a href="golang.org">go</a> installed first
<pre>
git clone https://github.com/Jisin0/Go-Filter-Bot
cd Go-Filter-Bot
go build .
./Go-Filter-Bot
</pre>
</p>
</details>

## Support

Ask any doubts or help in our support chat.
[![telegram badge](https://img.shields.io/badge/Telegram-Group-30302f?style=flat&logo=telegram)](https://telegram.dog/jisin_hub)

Join our telegram channel for more latest news and cool projects
[![telegram badge](https://img.shields.io/badge/Telegram-Channel-30302f?style=flat&logo=telegram)](https://telegram.dog/jisin_0)

## Thanks

 - Thanks To Paul For His Awesome [Library](https://github.com/PaulSonOfLars/gotgbot) And Support
 - Thanks To [Trojanz](https://github.com/trojanzhex) for Their Awesome [Unlimited Filter Bot](https://github.com/TroJanzHEX/Unlimited-Filter-Bot)

## Disclaimer
[![GNU General Public License 3.0](https://www.gnu.org/graphics/gplv3-127x51.png)](https://www.gnu.org/licenses/gpl-3.0.en.html#header)    
Licensed under [GNU GPL 3.0.](https://github.com/Jisin0/Go-Filter-Bot/blob/main/LICENSE).
Selling The Codes To Other People For Money Is *Strictly Prohibited*.
