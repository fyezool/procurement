import 'package:dio/dio.dart';
import '../models/user.dart';
import '../models/requisition.dart';
import '../models/purchase_order.dart';
import 'web_storage_service.dart';

class ApiService {
  final Dio _dio;

  // Private constructor
  ApiService._() : _dio = Dio(BaseOptions(baseUrl: 'http://localhost:8080/api')) {
    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final token = await WebStorageService.getToken();
          if (token != null) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          return handler.next(options);
        },
      ),
    );
  }

  // Singleton instance
  static final ApiService _instance = ApiService._();

  // Factory constructor to return the singleton instance
  factory ApiService() {
    return _instance;
  }

  // Public method to get the Dio instance
  Dio get dio => _dio;

  // Static methods for convenience
  static Future<Response<T>> get<T>(String path, {Map<String, dynamic>? queryParameters}) {
    return _instance.dio.get<T>(path, queryParameters: queryParameters);
  }

  static Future<Response<T>> post<T>(String path, dynamic data) {
    return _instance.dio.post<T>(path, data: data);
  }

  static Future<Response<T>> put<T>(String path, dynamic data) {
    return _instance.dio.put<T>(path, data: data);
  }

  static Future<Response<T>> delete<T>(String path) {
    return _instance.dio.delete<T>(path);
  }

  Future<List<User>> getUsers() async {
    try {
      final response = await _dio.get('/users');
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data;
        return data.map((userJson) => User.fromJson(userJson)).toList();
      } else {
        throw Exception('Failed to load users');
      }
    } catch (e) {
      throw Exception('Failed to load users: $e');
    }
  }

  Future<User> updateUser(int id, String name, String role) async {
    try {
      final response = await _dio.put(
        '/users/$id',
        data: {'name': name, 'role': role},
      );
      if (response.statusCode == 200) {
        return User.fromJson(response.data);
      } else {
        throw Exception('Failed to update user');
      }
    } catch (e) {
      throw Exception('Failed to update user: $e');
    }
  }

  Future<void> deleteUser(int id) async {
    try {
      final response = await _dio.delete('/users/$id');
      if (response.statusCode != 204) {
        throw Exception('Failed to delete user');
      }
    } catch (e) {
      throw Exception('Failed to delete user: $e');
    }
  }

  Future<List<Requisition>> getMyRequisitions() async {
    try {
      final response = await _dio.get('/requisitions/my');
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data;
        return data.map((json) => Requisition.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load requisitions');
      }
    } catch (e) {
      throw Exception('Failed to load requisitions: $e');
    }
  }

  Future<Requisition> updateRequisition(int id, Map<String, dynamic> data) async {
    try {
      final response = await _dio.put('/requisitions/$id', data: data);
      if (response.statusCode == 200) {
        return Requisition.fromJson(response.data);
      } else {
        throw Exception('Failed to update requisition');
      }
    } catch (e) {
      throw Exception('Failed to update requisition: $e');
    }
  }

  Future<void> deleteRequisition(int id) async {
    try {
      final response = await _dio.delete('/requisitions/$id');
      if (response.statusCode != 204) {
        throw Exception('Failed to delete requisition');
      }
    } catch (e) {
      throw Exception('Failed to delete requisition: $e');
    }
  }

  // Admin listing methods
  Future<List<Requisition>> getAllRequisitions() async {
    try {
      final response = await _dio.get('/requisitions/all');
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data;
        return data.map((json) => Requisition.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load all requisitions');
      }
    } catch (e) {
      throw Exception('Failed to load all requisitions: $e');
    }
  }

  Future<List<PurchaseOrder>> getAllPurchaseOrders() async {
    try {
      final response = await _dio.get('/purchase-orders/all');
      if (response.statusCode == 200) {
        final List<dynamic> data = response.data;
        return data.map((json) => PurchaseOrder.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load all purchase orders');
      }
    } catch (e) {
      throw Exception('Failed to load all purchase orders: $e');
    }
  }
}
