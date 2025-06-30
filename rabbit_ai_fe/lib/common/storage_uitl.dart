import 'package:shared_preferences/shared_preferences.dart';

///本地存储
class StorageUtils {
  StorageUtils._();

  /// 设置本地存储
  static Future<bool> setValue(String key, String value) async {
    final prefs = await SharedPreferences.getInstance();
    final status = await prefs.setString(key, value);
    return status;
  }

  /// 获取本地存储
  static Future<String> getValue(String key, [String defaultValue = '']) async {
    final prefs = await SharedPreferences.getInstance();
    final value = prefs.getString(key);
    return value ?? defaultValue;
  }

  /// 设置bool本地存储
  static Future<bool> setBoolValue(String key, bool value) async {
    final prefs = await SharedPreferences.getInstance();
    final status = await prefs.setBool(key, value);
    return status;
  }

  /// 获取bool本地存储
  static Future<bool> getBoolValue(String key,
      [bool defaultValue = false]) async {
    final prefs = await SharedPreferences.getInstance();
    final value = prefs.getBool(key);

    return value ?? defaultValue;
  }

  /// 移除本地存储(指定key)
  static Future<bool> remove(String key) async {
    final prefs = await SharedPreferences.getInstance();
    final status = await prefs.remove(key);

    return status;
  }

  /// 清除所有本地存储
  static Future<bool> clear() async {
    final prefs = await SharedPreferences.getInstance();
    final status = await prefs.clear();

    return status;
  }
}
