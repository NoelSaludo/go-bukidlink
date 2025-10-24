# Dart (Dio) usage — go-bukidlink

This file shows example usage of the go-bukidlink HTTP API using the Dart `dio` package. It provides:

- `pubspec.yaml` snippet to add `dio`
- Small model classes (`User`, `Item`, `Comment`) with `fromJson` / `toJson`
- Example API client using `Dio`
- Example calls for each route implemented by the server

> Notes:
> - The server runs at `http://localhost:8080` by default.
> - Replace sample UUIDs and values with those relevant to your environment.

---

## Add dio to your project

Add `dio` to your `pubspec.yaml` dependencies:

```yaml
dependencies:
  dio: ^5.0.0
  # add json serialization helpers if desired
```

Then run:

```bash
dart pub get
```

---

## Minimal models

These simple model classes are enough to parse server responses in examples below. In a real app you'd likely generate code with `json_serializable` or similar.

```dart
class User {
  final String id;
  final String username;
  final String password;
  final String email;

  User({required this.id, required this.username, required this.password, required this.email});

  factory User.fromJson(Map<String, dynamic> json) => User(
        id: json['id'] as String,
        username: json['username'] as String,
        password: json['password'] as String,
        email: json['email'] as String,
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'username': username,
        'password': password,
        'email': email,
      };
}

class Item {
  final String id;
  final String name;
  final String description;
  final int amount;
  final double costPKilo;
  final String category;
  final double? rating;

  Item({required this.id, required this.name, required this.description, required this.amount, required this.costPKilo, required this.category, this.rating});

  factory Item.fromJson(Map<String, dynamic> json) => Item(
        id: json['id'] as String,
        name: json['name'] as String,
        description: json['description'] as String,
        amount: (json['amount'] as num).toInt(),
        costPKilo: (json['costPKilo'] as num).toDouble(),
        category: json['category'] as String,
        rating: json['rating'] != null ? (json['rating'] as num).toDouble() : null,
      );
}

class Comment {
  final int id;
  final int itemid;
  final int userid;
  final String content;
  final double rating;

  Comment({required this.id, required this.itemid, required this.userid, required this.content, required this.rating});

  factory Comment.fromJson(Map<String, dynamic> json) => Comment(
        id: (json['id'] as num).toInt(),
        itemid: (json['itemid'] as num).toInt(),
        userid: (json['userid'] as num).toInt(),
        content: json['content'] as String,
        rating: (json['rating'] as num).toDouble(),
      );
}
```

---

## API client setup with Dio

```dart
import 'package:dio/dio.dart';

class ApiClient {
  final Dio _dio;

  ApiClient(String baseUrl)
      : _dio = Dio(BaseOptions(
          baseUrl: baseUrl,
          connectTimeout: const Duration(seconds: 5),
          receiveTimeout: const Duration(seconds: 5),
          headers: {'Accept': 'application/json'},
        ));

  // Helper to decode responses with error handling
  dynamic _unwrap(Response response) {
    if (response.statusCode != null && response.statusCode! >= 200 && response.statusCode! < 300) {
      return response.data;
    }
    throw DioException(requestOptions: response.requestOptions, response: response);
  }

  // Close client when done
  void close() => _dio.close();

  // Expose Dio for advanced uses
  Dio get dio => _dio;
}
```

Use like:

```dart
final api = ApiClient('http://localhost:8080');
```

---

## Examples

All examples assume `final api = ApiClient('http://localhost:8080');`.

### GET /ping

```dart
Future<void> ping(ApiClient api) async {
  try {
    final r = await api.dio.get('/ping');
    final data = r.data; // expected: {"message": "pong"}
    print('ping: $data');
  } on DioException catch (e) {
    print('ping failed: ${e.response?.statusCode} ${e.message}');
  }
}
```

### GET /item/:block

```dart
Future<List<Item>> getItemsBlock(ApiClient api, String block) async {
  final r = await api.dio.get('/item/$block');
  final data = r.data as List<dynamic>;
  return data.map((e) => Item.fromJson(e as Map<String, dynamic>)).toList();
}

// usage
// final items = await getItemsBlock(api, '0');
```

### GET /item/category/:category

```dart
Future<List<Item>> getItemsByCategory(ApiClient api, String category) async {
  try {
    final r = await api.dio.get('/item/category/$category');
    final list = (r.data as List<dynamic>).map((e) => Item.fromJson(e as Map<String, dynamic>)).toList();
    return list;
  } on DioException catch (e) {
    // server returns 500 for invalid category (per tests)
    print('error: ${e.response?.statusCode} ${e.message}');
    rethrow;
  }
}
```

### GET /user/:username

```dart
Future<User> getUser(ApiClient api, String username) async {
  final r = await api.dio.get('/user/$username');
  return User.fromJson(r.data as Map<String, dynamic>);
}

// usage
// final user = await getUser(api, 'JohnDoe');
```

### POST /user

```dart
Future<void> createUser(ApiClient api, User user) async {
  try {
    final r = await api.dio.post('/user', data: user.toJson(), options: Options(headers: {'Content-Type': 'application/json'}));
    if (r.statusCode == 201 || (r.statusCode != null && r.statusCode! >= 200 && r.statusCode! < 300)) {
      print('user created');
    } else {
      print('unexpected status: ${r.statusCode}');
    }
  } on DioException catch (e) {
    final status = e.response?.statusCode;
    if (status == 409) {
      print('user already exists (409)');
    } else {
      print('error creating user: $status ${e.message}');
    }
  }
}

// usage
// final newUser = User(id: '...', username: 'JohnDoe' ...);
// await createUser(api, newUser);
```

### GET /comment/:itemId

```dart
Future<List<Comment>> getComments(ApiClient api, String itemId) async {
  final r = await api.dio.get('/comment/$itemId');
  final list = (r.data as List<dynamic>).map((e) => Comment.fromJson(e as Map<String, dynamic>)).toList();
  return list;
}

// usage
// final comments = await getComments(api, 'a3e1b9f2-...');
```

### POST /comment/:productID (not implemented on server)

At the time of writing, `main.go` registers `commentG.POST('/:productID')` but no handler exists. Calling this endpoint will not create a comment. If/when the server implements a handler, you could POST JSON like this:

```dart
Future<void> postComment(ApiClient api, String productId, Map<String, dynamic> body) async {
  final r = await api.dio.post('/comment/$productId', data: body, options: Options(headers: {'Content-Type': 'application/json'}));
  print('status: ${r.statusCode}');
}

// example body (server dependent):
// {
//   "userid": 123,
//   "content": "Nice item",
//   "rating": 4.5
// }
```

---

## Tips & caveats

- If running on a real device/emulator, `localhost` might not point to your machine; use the machine IP or emulator host mapping (e.g., Android emulator `10.0.2.2`).
- For production, prefer configuring timeouts, retries, and proper error handling. Use `Interceptors` in `Dio` to add logging, retry, and auth headers.
- Consider adding typed API layers or using code generation for JSON models to avoid runtime casting issues.

---

## Example: small script

```dart
import 'dart:io';

void main() async {
  final api = ApiClient('http://localhost:8080');
  try {
    await ping(api);
    final items = await getItemsBlock(api, '0');
    print('fetched ${items.length} items');
  } finally {
    api.close();
  }
}
```

---

If you want, I can:

- Add a full Flutter example with state management.
- Provide generated models using `json_serializable`.
- Add an example that shows how to run the Dart script and test against a local server using Docker/Postgres.
