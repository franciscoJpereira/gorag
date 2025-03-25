# RAG API Server/TUI Application

This project, written in Golang, allows the user to have chats with a model by either running a server-side API or
running the system as a TUI application.
It also allows for Retrieval Augmented Generation (RAG) responses from the model by using a Chroma DB Knowledge base feed
with the data the user wants.
The main purpose of having a simple TUI application in this project is to provide a simple and interoperable interface to evaluate different models from different providers without having to build an interface for it.

-   The API is developed using the [Echo](https://echo.labstack.com) Golang framework. It exposes endpoints to create knwoledge bases, send either one-shot messages or full have full-on chats using the created knowledge bases.
- The TUI was built using [Bubbletea](https://github.com/charmbracelet/bubbletea) and exposes a simple interface to use the API without having to go through the hassle of building a Web or Desktop application to serve this purpose.

> The TUI has not yet the functionality to interact with Knowledge bases. This is still a work in progress.


## Running the TUI

To run the TUI application, you can either build the project as is or execute *go run main.go* while standing on the API
directory.

There is a **config.yaml** file in the *api/config* directory that allows to configure the models use for embedding and generation, as well as the endpoints to call to invocate these models.

## Running the project as an API

> Work in progress to configure how to run the whole system depending on an configuration parameter.

## Instantiating a Chroma DB fresh instance

You can instantiate a fresh Chroma DB instance by running the *clean-chroma* script that's provided in this repository.

To be able to run this script you must have chroma installed in your machine. You can do this by running:

```bash
    pip install chromadb
```

afterwards you can run the script by running it using python.

```bash
    python clean-chroma.py
```

The config file for the system points to this Chroma instance by default.

## TO-DOs

- :heavy_check_mark:  Possibility to run the system as either an API or TUI application based on config parameters.

- :x: Knowledge base support for the TUI application
    - :heavy_check_mark: Functionality to pick files
    - :x: Parse text files and send them to the backend
    - :x: Allow user to choose which KB to send a message to the LLM with

- :x: Chroma DB credentials support

- :x: Contenerization of the whole system 

- :x: Configuration to change the Model API interface that must be used.


