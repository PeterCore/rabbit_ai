import 'dart:io';

import 'package:android_id/android_id.dart';
import 'package:rabbit_ai_fe/common/extension_string.dart';
import 'package:device_info_plus/device_info_plus.dart';
import 'package:path_provider/path_provider.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:uuid/uuid.dart';

import 'app_platform_utils.dart';

class DeviceUtils {
  DeviceUtils._();

  static const _uuid = Uuid();

  static final _deviceInfo = DeviceInfoPlugin();

  static const _androidIdPlugin = AndroidId();

  /// 获取设备唯一ID
  static Future<String> getDeviceId() async {
    String id = "";
    if (AppPlatformUtils.isIOS) {
      id = await _getIOSId();
    } else if (AppPlatformUtils.isAndroid) {
      id = await _getAndroidId();
    } else if (AppPlatformUtils.isWeb) {
      //获取web的id
      WebBrowserInfo webInfo = await _deviceInfo.webBrowserInfo;
      id =
          "${webInfo.vendor ?? ""}${webInfo.userAgent ?? ""}${webInfo.hardwareConcurrency.toString()}"
              .toMd5();
    }
    return id;
  }

  /// 获取安卓唯一ID
  static Future<String> _getAndroidId() async {
    String? id = await _androidIdPlugin.getId();
    if (id == null || id.isEmpty) {
      id = await _getAndroidStorageId();
    }
    return id;
  }

  /// 获取IOS唯一ID
  static Future<String> _getIOSId() async {
    final iosInfo = await _deviceInfo.iosInfo;
    String id = iosInfo.identifierForVendor ?? '';
    if (id.isEmpty) {
      id = await _getIosStorageId();
    }
    return id;
  }

  // _deviceData = _readWebBrowserInfo(await _deviceInfo.webBrowserInfo);

  /// 获取IOS本地存储id
  static Future<String> _getIosStorageId() async {
    try {
      final file = await _getIosFile();
      final isNotExists = !file.existsSync();

      // 创建目录
      if (isNotExists) {
        file.createSync(recursive: true);
      }

      final storageId = await file.readAsString();

      if (storageId.isEmpty) {
        final String newId = _uuid.v1();

        await file.writeAsString(newId);

        return newId;
      } else {
        return storageId;
      }
    } catch (_) {
      return '';
    }
  }

  /// 获取IOS文件路径
  static Future<File> _getIosFile() async {
    final directory = await getApplicationDocumentsDirectory();
    return File('${directory.path}/.uuid.txt');
  }

  /// 获取安卓文件路径
  static Future<File> _getAndroidFile() async {
    final androidInfo = await _deviceInfo.androidInfo;
    final version = double.parse(androidInfo.version.release);
    final isExternal = version >= 11;

    final bool isStorageGranted = await Permission.storage.request().isGranted;
    final bool isExternalGranted =
        await Permission.manageExternalStorage.request().isGranted;

    // 低版本授权
    final bool isLow = !isExternal && isStorageGranted;

    // 高版本授权（11及以上）
    final bool isHigh = isExternal && isStorageGranted && isExternalGranted;

    if (isLow || isHigh) {
      return File('/storage/emulated/0/.flutter/.uuid.txt');
    } else {
      final directory = await getApplicationDocumentsDirectory();
      return File('${directory.path}/.uuid.txt');
    }
  }

  /// 获取安卓本地存储id
  static Future<String> _getAndroidStorageId() async {
    try {
      final file = await _getAndroidFile();
      final isNotExists = !file.existsSync();

      // 创建目录
      if (isNotExists) {
        file.createSync(recursive: true);
      }

      final storageId = await file.readAsString();

      if (storageId.isEmpty) {
        final String newId = _uuid.v1();

        await file.writeAsString(newId);

        return newId;
      } else {
        return storageId;
      }
    } catch (_) {
      return '';
    }
  }
}
