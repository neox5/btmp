// Package btmp provides a compact, growable bitmap optimized for fast range
// updates and overlap-safe copies, plus a zero-copy row-major 2D Grid view.
//
// Conventions:
//   - Length is in bits (Len).
//   - Storage is []uint64 words, exposed via Words() for read-only inspection.
//   - Ranges use (start, count).
//   - SetRange and CopyRange auto-grow; ClearRange is in-bounds.
//   - All mutating methods return *Bitmap for chaining.
//
// Invariant:
//   - After any public mutator returns, all bits at indexes >= Len() are zero,
//     even when count == 0.
package btmp
