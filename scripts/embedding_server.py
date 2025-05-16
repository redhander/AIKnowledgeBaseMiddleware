from sentence_transformers import SentenceTransformer
from flask import Flask, request, jsonify

app = Flask(__name__)
model = SentenceTransformer('all-MiniLM-L6-v2')  # 也可以换成 'text-embedding-3-small'

@app.route('/embed', methods=['POST'])
def embed():
    texts = request.json.get('texts')
    vectors = model.encode(texts)
    return jsonify({"vectors": vectors.tolist()})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)