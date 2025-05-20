from sentence_transformers import SentenceTransformer
from flask import Flask, request, jsonify

app = Flask(__name__)
print("Loading model all-MiniLM-L6-v2...")
# model = SentenceTransformer('all-MiniLM-L6-v2')  # 也可以换成 'text-embedding-3-small'
model = SentenceTransformer('/models/all-MiniLM-L6-v2')  # 改为从容器内路径加载模型

print("Model loaded successfully.")
#model.save('./all-MiniLM-L6-v2')
@app.route('/embed', methods=['POST'])
def embed():
    data = request.get_json()  # 获取完整的 JSON 数据
    if not data or 'texts' not in data:
        return jsonify({"error": "Missing 'texts' field in request body"}), 400

    texts = data['texts']
    if not isinstance(texts, list):
        return jsonify({"error": "'texts' must be a list"}), 400

    vectors = model.encode(texts)
    return jsonify({"vectors": vectors.tolist()})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)