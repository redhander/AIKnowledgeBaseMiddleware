FROM python:3.9-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install huggingface_hub[hf_xet]
COPY models /models
COPY embedding_server.py .
EXPOSE 8081

CMD ["python", "embedding_server.py"]
