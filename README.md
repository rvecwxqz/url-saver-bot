# url-saver-bot

URL saver is a Telegram bot written in Golang that utilizes a multilingual BERT model for natural language processing. The machine learning part of the bot is implemented in Python using PyTorch GPU. Communication between the Go and the Python BERT model is done via gRPC.

## Features

- Save and tag links: Send a link to the bot and it will automatically extract the title and content of the webpage, and suggest tags based on the content using the BERT model.
- Easy retrieval: Use tags to quickly search for and filter your saved links.
- Command-driven interface: Interact with the bot using commands such as /get, /show_tags, /show_all, and /remove.

## Prerequisites

- Golang
- Python
- PyTorch GPU
- Transformers
- PostgreSQL database
- Telegram Bot API token
- Pre-trained multilingual BERT [Google drive](https://drive.google.com/file/d/1kTJC3X9RTHXeiqoEPRJi_CfG7UiXiODX/view?usp=sharing). Put in /internal/ml/bert-classifier/model/content/


## Build and run the bot:

1. Clone this repository:

    git clone https://github.com/rvecwxqz/url-saver-bot.git

2. Download the pre-trained multilingual BERT model and place it in the `/internal/ml/bert-classifier/model/content/` directory.

3. Put `'YOUR_TELEGRAM_API_TOKEN'` in the file `/internal/config/config.go`.

4. Build and run the bot using the following command in the project's root directory:

    go run cmd/app/main.go

5. Start a conversation with your bot on Telegram and use the available commands to save, retrieve, and manage your links.
