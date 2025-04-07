FROM python:latest
WORKDIR /backend
EXPOSE 80
COPY ./requirements.txt ./
RUN pip install --no-cache-dir --upgrade pip
RUN pip install --no-cache-dir -r ./requirements.txt
COPY . ./
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]
