import 'dart:html' as html;
import 'dart:convert';

class WebStorageService {
  static const String _tokenKey = 'auth_token';
  static const String _userKey = 'user_data';

  // Store JWT token
  static Future<void> storeToken(String token) async {
    html.window.localStorage[_tokenKey] = token;
  }

  // Get JWT token
  static Future<String?> getToken() async {
    return html.window.localStorage[_tokenKey];
  }

  // Store user data
  static Future<void> storeUserData(Map<String, dynamic> userData) async {
    html.window.localStorage[_userKey] = jsonEncode(userData);
  }

  // Get user data
  static Future<Map<String, dynamic>?> getUserData() async {
    final userData = html.window.localStorage[_userKey];
    if (userData != null) {
      return jsonDecode(userData) as Map<String, dynamic>;
    }
    return null;
  }

  // Clear all stored data
  static Future<void> clearAll() async {
    html.window.localStorage.remove(_tokenKey);
    html.window.localStorage.remove(_userKey);
  }
}
