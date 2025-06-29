import 'package:intl/intl.dart';

class AppUtils {
  static DateFormat dateFormat = DateFormat("yyyy-MM-dd");
  static DateFormat dateTimeFormat = DateFormat("MM-dd HH:mm");
  static DateFormat dateTimeFormatWithYear = DateFormat("yyyy-MM-dd HH:mm");
  static DateFormat dateTimeFormatWithYearSecond =
      DateFormat("yyyy-MM-dd HH:mm:ss");
  static DateFormat dateTime12formatter = DateFormat("MM-dd a h:mm", "zh_CN");

  static String formatFriendlyTime(int timestampMillis) {
    final now = DateTime.now();
    final dt = DateTime.fromMillisecondsSinceEpoch(timestampMillis);

    // 只保留年月日，用来算相差天数
    final today = DateTime(now.year, now.month, now.day);
    final thatDay = DateTime(dt.year, dt.month, dt.day);

    final daysDiff = today.difference(thatDay).inDays;

    if (daysDiff == 0) {
      // 同一天：刚刚 / X 分钟前 / X 小时前
      final diff = now.difference(dt);
      if (diff.inMinutes < 1) {
        return '刚刚';
      } else if (diff.inHours < 1) {
        return '${diff.inMinutes}分钟前';
      } else {
        return '${diff.inHours}小时前';
      }
    } else if (daysDiff == 1) {
      // 昨天
      return '昨天';
    } else {
      // 不是昨天、也不是今天
      if (dt.year == now.year) {
        // 同一年：MM-dd
        return DateFormat('MM-dd').format(dt);
      } else {
        // 往年：yyyy-MM-dd
        return DateFormat('yyyy-MM-dd').format(dt);
      }
    }
  }

  /// 时间戳格式化-秒
  static String formatTimestamp(int ts) {
    if (ts == 0) {
      return "----";
    }
    return formatTimestampMS(ts * 1000);
  }

  static String format12TimestampMS(int ts) {
    if (ts == 0) {
      return "----";
    }
    DateTime dt = DateTime.fromMillisecondsSinceEpoch(ts);
    return dateTime12formatter.format(dt);
  }

  static String formatYearMonthTimestamp(int ts) {
    if (ts == 0) {
      return "----";
    }
    return formatYearMonthTimestampMS(ts * 1000);
  }

  static String formatTimestampToDate(int ts) {
    if (ts == 0) {
      return "----";
    }
    var dt = DateTime.fromMillisecondsSinceEpoch(ts * 1000);
    return dateFormat.format(dt);
  }

  static String formatYearMonthTimestampMS(int ts) {
    var dt = DateTime.fromMillisecondsSinceEpoch(ts * 1000);
    return "${dt.year.toString()}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}";
  }

  static String formatYearMonthSTimestampMS(int ts) {
    var dt = DateTime.fromMillisecondsSinceEpoch(ts * 1000);
    return dateTimeFormatWithYearSecond.format(dt);
  }

  /// 时间戳格式化-毫秒
  static String formatTimestampMS(int ts) {
    var dt = DateTime.fromMillisecondsSinceEpoch(ts);

    var dtNow = DateTime.now();
    if (dt.year == dtNow.year &&
        dt.month == dtNow.month &&
        dt.day == dtNow.day) {
      return "今天";
    }
    if (dt.year == dtNow.year &&
        dt.month == dtNow.month &&
        dt.day == dtNow.day - 1) {
      return "昨天";
    }
    if (dt.day != dtNow.day) {
      return "${dt.month.toString().padLeft(2, '0')}月${dt.day.toString().padLeft(2, '0')}日 ${dt.year.toString()}年";
    }

    return dateTimeFormatWithYear.format(dt);
  }
}
