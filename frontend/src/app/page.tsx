import Link from 'next/link'

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-900 to-blue-700 text-white">
      <div className="container mx-auto px-4 py-16">
        <header className="flex justify-between items-center mb-16">
          <h1 className="text-2xl font-bold">Church Management System</h1>
          <nav>
            <Link 
              href="/login" 
              className="px-4 py-2 bg-white text-blue-900 rounded-lg hover:bg-blue-100 transition-colors"
            >
              Login
            </Link>
          </nav>
        </header>
        
        <main className="max-w-4xl mx-auto text-center">
          <h2 className="text-5xl font-bold mb-6">Fruit Management System</h2>
          <p className="text-xl mb-12 text-blue-100">
            Track your ministry's growth with our comprehensive contact management system.
            Monitor Bible studies, spiritual progression, and manage room reservations all in one place.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link 
              href="/login" 
              className="px-8 py-4 bg-white text-blue-900 font-bold rounded-lg hover:bg-blue-100 transition-colors"
            >
              Get Started
            </Link>
            <Link 
              href="/register" 
              className="px-8 py-4 border-2 border-white rounded-lg hover:bg-blue-800 transition-colors"
            >
              Create Account
            </Link>
          </div>
          
          <div className="mt-24 grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="bg-blue-800 p-6 rounded-lg">
              <div className="text-3xl mb-4">👥</div>
              <h3 className="text-xl font-bold mb-2">Contact Management</h3>
              <p>Track spiritual growth from first contact to gospel worker.</p>
            </div>
            
            <div className="bg-blue-800 p-6 rounded-lg">
              <div className="text-3xl mb-4">📚</div>
              <h3 className="text-xl font-bold mb-2">Bible Studies</h3>
              <p>Record and monitor the 30 studies needed for gospel worker status.</p>
            </div>
            
            <div className="bg-blue-800 p-6 rounded-lg">
              <div className="text-3xl mb-4">🏢</div>
              <h3 className="text-xl font-bold mb-2">Room Reservations</h3>
              <p>Easily book and manage spaces for meetings and studies.</p>
            </div>
          </div>
        </main>
        
        <footer className="mt-24 text-center text-blue-200">
          <p>&copy; {new Date().getFullYear()} Church Management System</p>
        </footer>
      </div>
    </div>
  )
}
