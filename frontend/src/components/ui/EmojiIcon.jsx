export default function EmojiIcon({ children, dark = false, size = 'md' }) {
  const sizeClasses = {
    sm: 'text-base',
    md: 'text-xl',
    lg: 'text-2xl',
    xl: 'text-3xl'
  };

  // Use dark skin tone variants for emojis when dark=true (African context)
  const darkSkinEmojis = {
    '👨‍⚕️': '👨🏿‍⚕️',  // Health worker: dark skin tone
    '👥': '👥',        // People silhouette (no skin tone variants)
    '💰': '💰',        // Money bag (no skin tone)
  };

  const displayEmoji = dark && darkSkinEmojis[children]
    ? darkSkinEmojis[children]
    : children;

  return (
    <span className={`inline-block ${sizeClasses[size]}`}>
      {displayEmoji}
    </span>
  );
}
