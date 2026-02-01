set -e

echo "🔄 Running migrations..."
/app/migrator

echo "✅ Migrations completed successfully"

echo "🚀 Starting application..."
exec /app/app
