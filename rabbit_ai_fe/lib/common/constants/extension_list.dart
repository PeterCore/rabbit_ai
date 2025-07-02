extension ExtensionList<T> on List<T> {
  List<T> takeLengthThree(int subLength) {
    return length <= subLength ? this : sublist(0, subLength);
  }

  List<T> rotatedLeft() {
    if (isEmpty) return this;
    return sublist(1) + [first];
  }
}
