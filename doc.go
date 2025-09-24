// btmp provides a compact, growable bitmap optimized for fast range updates
// and overlap-safe bit copying. Includes a zero-copy row-major Grid view for 2D use.
//
// Conventions:
//   - Length is in bits (Len).
//   - Storage is []uint64 words.
//   - Ranges use (start, count).
//   - SetRange and CopyRange auto-grow; ClearRange is in-bounds.
//   - All mutating methods return *Bitmap for chaining (pointer never changes).
//   - Grid constructors return *Grid; Grid mutators return *Grid for chaining.
package btmp
