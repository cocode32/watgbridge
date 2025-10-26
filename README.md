# ⚠️ Disclaimer
This project is provided as-is and without any warranty or guarantee of any kind.
Use it entirely at your own risk.

The maintainers and contributors are not responsible for any damages, losses, account suspensions, bans, or other consequences that may arise from the use of this software.

Some platforms or services (such as WhatsApp or similar messaging systems) may restrict or prohibit automated or third-party interactions. You are solely responsible for ensuring that your use of this project complies with all applicable terms of service and laws.
By using or distributing this software, you acknowledge that:
- You understand the potential risks involved.
- You accept full responsibility for any actions or outcomes resulting from its use.

# coco_sanity changes

I've kept the original README.md intact below. I'm adding additional information here for anyone that would like to use this instead of the original.

I don't currently have time to make this backwards compatible for all the users of the original `watgbridge`. So, if someone wants to cherry-pick commits from this branch, to merge it into the original; please go ahead.

> ⚠️
>
> DO NOT USE THIS ON EXISTING DATABASE
>
> The database has been completely redone as a result of my attempt to gain sanity back in my life.

## Database Changes
All the `"new"` database tables are prefixed with `Coco`. Just so that I can keep sanity

# License Changes

This project is a fork of [watgbridge](https://github.com/akshettrj/watgbridge), originally licensed under the MIT License.

All modifications, additions, and new code introduced in this fork are licensed under the **GNU Affero General Public License v3.0 (AGPLv3)**.  
The original code from [Original Project Name] remains under its **MIT License**.

### Summary of Licensing
- **Original Code:** MIT License (© original authors)
- **Enhancements & Additions:** GNU AGPLv3 (© you, [Newt Fourie])
- **Combined Distribution:** As a whole, this fork is distributed under the terms of the AGPLv3.

See the [LICENSE](./LICENSE) file for full details.


# Release Notes

## 2025 October

The main change that I've implemented here is getting a "somewhat" consistent mapping of JID and LID for users. I don't have enough test data to confirm, but from what I can see, it should be good.

> ⚠️
>
> VERY IMPORTANT
> If you decide to use this head, then you need to make sure that you completely delete your current database, and the wawebstore database. Both these database need to be removed, because I'm not using the same ID mapping that was used from the original code. You can leave the config as is, but the database will need to be restarted. You'll also basically get new threads for each chat, because nothing will match, so you might want to create a new Telegram group.

I would love to make this change, and in place upgrade, but from what I've read about the LIDs here: [https://support.whapi.cloud/help-desk/groups/what-is-lid-in-whatsapp-groups](https://support.whapi.cloud/help-desk/groups/what-is-lid-in-whatsapp-groups) it actually seems like they're (whatsapp) are not done with these changes. I'm just getting super annoyed with lids and phone numbers not matching the same chat. So, that's the purpose of this new head

# Installation & Setup

For ease, I've updated the installation instructions too, because it's just easier for me when setting up a new instance to follow these setup instructions. They are the same as what's at the bottom, just modified for easier use.

basically just follow the below commands in the terminal
```bash
sudo adduser --group --home /opt/watgbridge watgbridge
sudo passwd -l watgbridge
sudo apt install git gcc golang ffmpeg imagemagick -y
sudo su
# login as root to go into the directory and build the source
cd /opt/watgbridge
# make sure to clone the repo directly into this folder -- note the ./ at the end
sudo -u watgbridge git clone https://github.com/cocode32/watgbridge.git ./
sudo -u watgbridge go build
# MAKE ALL THE CONFIG CHANGES
#### THEN ####
# start the app manually to get the QR code
sudo -u watgbridge ./watgbridge
```

---
# Original README
---
# WhatsApp-Telegram-Bridge

Despite the name, it's not exactly a "bridge". It forwards messages from WhatsApp to Telegram and you can reply to them
from Telegram.

<a href="https://t.me/PropheCProjects">
  <img src="https://img.shields.io/badge/Original_Author_Updates_Channel-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white"></img>
</a>&nbsp; &nbsp;
<a href="https://t.me/WaTgBridge">
  <img src="https://img.shields.io/badge/Original_Author_Discussion_Group-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white"></img>
</a>&nbsp; &nbsp;
<a href="https://youtu.be/xc75XLoTmA4">
  <img src="https://img.shields.io/badge/Setup_YouTube_Video-FF0000?style=for-the-badge&logo=youtube&logoColor=white"</img>
</a>

# DISCLAIMER !!!

This project is in no way affiliated with WhatsApp or Telegram. Using this can also lead to your account getting banned by WhatsApp so use at your own risk.

## Sample Screenshots

<p align="center">
  <img src="./assets/telegram_side_sample.png" width="350" alt="Telegram Side">
  <img src="./assets/whatsapp_side_sample.jpg" width="350" alt="WhatsApp Side">
</p>

## Features and Design Choices

- All messages from various chats (on WhatsApp) are sent to different topics/threads within the same target group (on Telegram)
- Configuration options available to disable different types of updates from WhatsApp
- Can reply and send new messages from Telegram
- Can tag all people using @all or @everyone. Others can also use this in group chats which you specify in configuration file
- Can react to messages by replying with single instance of the desired emoji
- Supports static stickers from both ends
- Can send Animated (TGS) stickers from Telegram
- Video stickers from Telegram side are supported
- Video stickers from WhatsApp side are currently forwarded as GIFs to Telegram

## Bugs and TODO

- Document naming is messed up and not consistent on Telegram, have to find a way to always send same names

PRs are welcome :)


## Installation

- Make a supergroup (enable message history for new members) with topics enabled
- Add your bot in the group, make it an admin with permissions to `Manage topics`
- Install `git`, `gcc` and `golang`, `ffmpeg` , `imagemagick` (optional), on your system
- Clone this repository anywhere and navigate to the cloned directory
- Run `go build`
- Copy `sample_config.yaml` to `config.yaml` and fill the values, there are comments to help you.
- Execute the binary by running `./watgbridge`
- On first run, it will show QR code for logging into WhatsApp that can be~~~~ scanned by the WhatsApp app in `Linked devices`
- It is recommended to restart the bot after every few hours becuase WhatsApp likes to disconnect a lot. So a sample Systemd service file has been provided (`watgbridge.service.sample`). Edit the `User` and `ExecStart` according to your setup:
    - If you do not have local bot API server, remove `tgbotapi.service` from the `After` key in `Unit` section.
    - This service file will restart the bot every 24 hours
