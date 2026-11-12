// This file implements assigning and releasing resources. These operations will run atomically
// so concurrent dispatch requests cannot assign the same resource more than once.
package application
