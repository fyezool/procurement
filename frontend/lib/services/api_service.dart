import 'package:dio/dio.dart';
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
}
