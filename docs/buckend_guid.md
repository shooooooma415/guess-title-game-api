# Backend Implementation Guide

## 必須データ構造

### ゲームデータ保持
```javascript
let gameData = {
  topic: null,
  originalEmojis: [],    // ホストが選んだ絵文字(3~5)
  displayedEmojis: [],   // ダミー込み4~6つ
  dummyIndex: null,      // 0-(3-5)
  dummyEmoji: null,      // "🎭"
  answer: null,
  assignments: []
};
```

### 参加者データ
```javascript
participants = [
  {
    user_id: "id",
    user_name: "name",
    role: "host" | "player",
    is_Leader: true | false
  }
];
```

**重要:** 最初の参加者を `is_Leader: true` に設定

---

## タイマー設定

- **議論時間:** 5分（300秒）
- **開始遅延:** 5秒
- **フォーマット:** "MM:SS"
- **送信頻度:** 毎秒（TIMER_TICK）

---

## 状態遷移

```
WAITING → SETTING_TOPIC → DISCUSSING → ANSWERING → CHECKING → FINISHED
```

---

## HTTP API

### POST /api/rooms
```json
Response: {
  "room_id": "abc123",
  "user_id": "host-id",
  "room_code": "AAAAAA",
  "theme": "人物",
  "hint": "hint"
}
```

### POST /api/user
```json
Request: { "room_code": "AAAAAA", "user_name": "name" }
Response: { "room_id": "abc123", "user_id": "id", "is_leader": true }
```
→ 最初の参加者を `is_Leader: true` に設定

### POST /api/rooms/:room_id/start
権限: `role === "host"`  
→ STATE_UPDATE (setting_topic) 送信

### POST /api/rooms/:room_id/topic
権限: `role === "host"`  
```json
Request: { "topic": "topic", "emojis": ["🍎", "📱", "👔"] }
```
→ WebSocket SUBMIT_TOPIC を待ってマージ → DISCUSSING へ → 5秒後タイマー開始

### POST /api/rooms/:room_id/answer
権限: `is_Leader === true`  
```json
Request: { "user_id": "id", "answer": "answer" }
```
→ CHECKING へ

### POST /api/rooms/:room_id/skip-discussion
権限: `role === "host"`  
→ タイマークリア → ANSWERING へ（**ダミーデータ必須**）

### POST /api/rooms/:room_id/finish
権限: `role === "host"`  
→ FINISHED へ

---

## WebSocket メッセージ

### クライアント → サーバー

**CLIENT_CONNECTED**
```json
{ "type": "CLIENT_CONNECTED", "payload": { "user_id": "id", "user_name": "name" } }
```
→ PARTICIPANT_UPDATE 配信

**FETCH_PARTICIPANTS**
```json
{ "type": "FETCH_PARTICIPANTS" }
```
→ 参加者リスト返送

**SUBMIT_TOPIC**
```json
{
  "type": "SUBMIT_TOPIC",
  "payload": {
    "displayedEmojis": ["🍎", "📱", "👔", "🎭"],
    "originalEmojis": ["🍎", "📱", "👔"],
    "dummyIndex": 3,
    "dummyEmoji": "🎭"
  }
}
```
→ HTTP /topic の後に送信 → ダミーデータ保存 → DISCUSSING へ → 5秒後タイマー

**ANSWERING**
```json
{
  "type": "ANSWERING",
  "payload": {
    "answer": "answer",
    "displayedEmojis": [...],
    "originalEmojis": [...],
    "dummyIndex": 3,
    "dummyEmoji": "🎭"
  }
}
```
→ 全データ保存 → CHECKING へ

---

### サーバー → クライアント

**STATE_UPDATE**
```json
{
  "type": "STATE_UPDATE",
  "payload": {
    "nextState": "discussing",
    "data": {
      "topic": "topic",
      "displayedEmojis": [...],
      "originalEmojis": [...],
      "dummyIndex": 3,
      "dummyEmoji": "🎭",
      "assignments": [...]
    }
  }
}
```
nextState: `setting_topic` | `discussing` | `answering` | `checking` | `finished`

**PARTICIPANT_UPDATE**
```json
{ "type": "PARTICIPANT_UPDATE", "payload": { "participants": [...] } }
```

**TIMER_TICK**
```json
{ "type": "TIMER_TICK", "payload": { "time": "04:59" } }
```

**ERROR**
```json
{ "type": "ERROR", "payload": { "code": "code", "message": "msg" } }
```

---

## 重要な実装ポイント

### ✅ ダミーデータは全状態遷移で送信
- DISCUSSING → ANSWERING → CHECKING で `displayedEmojis`, `originalEmojis`, `dummyIndex`, `dummyEmoji` を**必ず含める**
- skip-discussion 時も**必須**

### ✅ HTTP + WebSocket 二重通信
```javascript
// HTTP /topic で topic 保存
POST /topic -> gameData.topic = body.topic;

// WebSocket SUBMIT_TOPIC でダミーデータ保存
WS SUBMIT_TOPIC -> {
  gameData.displayedEmojis = payload.displayedEmojis;
  gameData.originalEmojis = payload.originalEmojis;
  gameData.dummyIndex = payload.dummyIndex;
  gameData.dummyEmoji = payload.dummyEmoji;
  broadcast(STATE_UPDATE); // 全データ送信
}
```

### ✅ タイマーは5秒後に開始
```javascript
broadcast(STATE_UPDATE); // DISCUSSING へ
setTimeout(() => {
  // 5秒待ってからタイマー開始（300秒）
  setInterval(...);
}, 5000);
```

### ✅ 権限チェック
- ホスト操作: `role === "host"`
- リーダー操作: `is_Leader === true`  
  ※ `role` と `is_Leader` は別の概念