# Setup Guide: https://docs.maltego.com/support/solutions/articles/15000034047
FROM ubuntu:18.04

# Install system and Python dependencies
RUN apt update
RUN apt install python3 python3-pip -y

# Configure Python3 and pip as the default for your system
RUN echo "alias python=python3" >> ~/.bash_aliases
RUN echo "alias pip=pip3" >> ~/.bash_aliases
RUN . ~/.bash_aliases

# Install Maltego Python Transform Library
RUN pip3 install maltego-trx

# Install Gunicorn:
RUN pip3 install gunicorn

RUN mkdir -p /var/www/
COPY ./TRX /var/www/TRX
RUN chown -R www-data:www-data /var/www/TRX/

WORKDIR /var/www/TRX/
ENTRYPOINT ["gunicorn", "--bind=0.0.0.0:8081", "--threads=25", "--workers=2", "project:app"]