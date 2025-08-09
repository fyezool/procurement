import 'package:flutter/foundation.dart';
import 'web_storage_service.dart';
import 'api_service.dart';
import 'package:dio/dio.dart';
import 'dart:convert';

class AuthService extends ChangeNotifier {
  String? _token;
  Map<String, dynamic>? _user;
  bool _isAuthenticated = false;

  bool get isAuthenticated => _isAuthenticated;
  String? get token => _token;
  Map<String, dynamic>? get user => _user;
  String? get userRole => _user != null ? _user!['role'] : null;

  AuthService() {
    _loadStoredAuth();
  }

  Future<void> _loadStoredAuth() async {
    _token = await WebStorageService.getToken();
    if (_token != null) {
      _user = _decodeToken(_token!);
      _isAuthenticated = true;
    } else {
      _isAuthenticated = false;
    }
    notifyListeners();
  }

  Future<bool> login(String email, String password) async {
    try {
      final response = await ApiService.post('/login', {
        'email': email,
        'password': password,
      });

      if (response.statusCode == 200) {
        _token = response.data['token'];
        _user = _decodeToken(_token!);
        _isAuthenticated = true;

        await WebStorageService.storeToken(_token!);

        notifyListeners();
        return true;
      }
      return false;
    } on DioError catch (e) {
      print('Login error: $e');
      return false;
    }
  }

  Future<void> logout() async {
    _token = null;
    _user = null;
    _isAuthenticated = false;

    await WebStorageService.clearAll();
    notifyListeners();
  }

  Map<String, dynamic> _decodeToken(String token) {
    try {
      final parts = token.split('.');
      if (parts.length != 3) {
        throw Exception('Invalid token');
      }
      final payload = parts[1];
      final normalized = base64Url.normalize(payload);
      final resp = utf8.decode(base64Url.decode(normalized));
      return json.decode(resp);
    } catch (e) {
      print('Error decoding token: $e');
      return {};
    }
  }
}
