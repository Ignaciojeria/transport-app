import Link from 'next/link'
import Image from 'next/image'

export const metadata = {
  title: 'Terms and Conditions - MiCartaPro',
  description: 'Terms and conditions of use for MiCartaPro. Learn about the terms governing the use of our services.',
}

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
      {/* Navigation */}
      <nav className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link href="/" className="flex items-center space-x-2">
              <Image 
                src="/logo.png" 
                alt="MiCartaPro Logo" 
                width={240} 
                height={72}
                className="h-16 md:h-20 w-auto"
              />
            </Link>
            <Link 
              href="/"
              className="text-gray-600 hover:text-blue-600 transition-colors"
            >
              Back to Home
            </Link>
          </div>
        </div>
      </nav>

      {/* Content */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-16">
        <div className="bg-white rounded-lg shadow-lg p-8 md:p-12">
          <h1 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            Terms and Conditions of Use
          </h1>
          
          <p className="text-gray-600 mb-8">
            <strong>Last updated:</strong> {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>

          <div className="prose prose-lg max-w-none space-y-8 text-gray-700">
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. Acceptance of Terms</h2>
              <p>
                Welcome to MiCartaPro. These Terms and Conditions of Use ("Terms") govern your access and use of our digital menu services, including our website, applications, and any related services (collectively, the "Services").
              </p>
              <p>
                By accessing or using our Services, you agree to be bound by these Terms. If you do not agree to any part of these Terms, you must not use our Services.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Description of Services</h2>
              <p>
                MiCartaPro is a self-service SaaS (Software as a Service) platform that allows restaurants to create, manage, and share digital menus with their customers. Our platform provides:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Self-service tools for customizing digital menu designs</li>
                <li>QR code generation functionality</li>
                <li>Integrated shopping cart system</li>
                <li>WhatsApp integration for order reception</li>
                <li>Responsive design for all devices</li>
                <li>Continuous platform improvements, new features, and updates</li>
              </ul>
              <p className="mt-4">
                <strong>Important:</strong> MiCartaPro is a self-service platform. We do not provide ongoing human-driven services such as custom design work, manual menu creation, or personalized consulting as part of the subscription. All plans include access to the platform and ongoing software improvements.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Eligibility</h2>
              <p>
                To use our Services, you must:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Be at least 18 years of age or have the consent of a parent or legal guardian</li>
                <li>Have the legal capacity to enter into binding contracts</li>
                <li>Not be prohibited from using our Services under applicable laws</li>
                <li>Provide accurate and complete information when registering or using our Services</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Use of Services</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.1. Permitted Use</h3>
              <p>You may use our Services solely for lawful purposes and in accordance with these Terms. You agree to:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Use the Services responsibly and ethically</li>
                <li>Respect the intellectual property rights of others</li>
                <li>Comply with all applicable laws and regulations</li>
                <li>Not interfere with the operation of the Services</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.2. Prohibited Use</h3>
              <p>The following is strictly prohibited:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Using the Services for illegal or fraudulent activities</li>
                <li>Attempting unauthorized access to our systems or data</li>
                <li>Transmitting viruses, malware, or malicious code</li>
                <li>Reverse engineering, decompiling, or disassembling our Services</li>
                <li>Using bots, automated scripts, or similar methods to access the Services</li>
                <li>Copying, modifying, or distributing our Services without authorization</li>
                <li>Using the Services to send spam or unsolicited communications</li>
                <li>Violating the privacy or intellectual property rights of others</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Accounts and Registration</h2>
              <p>
                To access certain features of our Services, you may need to create an account. You are responsible for:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Maintaining the confidentiality of your account credentials</li>
                <li>All activities that occur under your account</li>
                <li>Notifying us immediately of any unauthorized use of your account</li>
                <li>Providing accurate and up-to-date information</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. User Content</h2>
              <p>
                When using our Services, you may provide content, including menu information, images, text, and other materials ("User Content"). By providing User Content, you grant MiCartaPro a non-exclusive, worldwide, royalty-free, and transferable license to use, reproduce, modify, and distribute such content solely for the purpose of providing and improving our Services.
              </p>
              <p>
                You are responsible for ensuring that your User Content:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Does not infringe the intellectual property rights of third parties</li>
                <li>Does not contain defamatory, obscene, or illegal material</li>
                <li>Does not violate the privacy rights of others</li>
                <li>Complies with all applicable laws and regulations</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. Intellectual Property</h2>
              <p>
                All rights, titles, and interests in and to our Services, including but not limited to software, design, text, graphics, logos, icons, and data compilations, are the property of MiCartaPro or its licensors and are protected by international intellectual property laws.
              </p>
              <p>
                You are not granted any right, title, or interest in our Services except the limited right to use in accordance with these Terms.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Pricing and Payments</h2>
              <p>
                Our services are available through a simple monthly subscription:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Monthly Subscription:</strong> $15 USD per month</li>
                <li>No setup fees or long-term commitments</li>
                <li>Cancel anytime</li>
                <li>Prices are subject to change with prior notice</li>
                <li>Payments are processed by Paddle, who acts as the Merchant of Record</li>
              </ul>
              <p className="mt-4">
                All prices are expressed in US dollars (USD) and may be subject to taxes according to your jurisdiction.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8.1. Refund Policy</h2>
              <p>
                We offer a <strong>14-day money-back guarantee</strong> on all purchases made via Paddle.
              </p>
              <p>
                If you're not satisfied with your purchase, you can request a full refund within 14 days from the date of purchase by contacting Paddle or our support team.
              </p>
              <p>
                After 14 days, refunds are handled at our and Paddle's discretion.
              </p>
              <p className="mt-4">
                For detailed information about our refund policy, please refer to our <Link href="/refund" className="text-blue-600 hover:text-blue-700 underline">Refund Policy</Link> page.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">9. Service Availability</h2>
              <p>
                We strive to keep our Services available continuously, but we do not guarantee that the Services will be available at all times, free from interruptions or errors. We may perform scheduled or unscheduled maintenance that may result in temporary service interruptions.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">10. Limitation of Liability</h2>
              <p>
                TO THE MAXIMUM EXTENT PERMITTED BY APPLICABLE LAW, MICARTAPRO AND ITS AFFILIATES SHALL NOT BE LIABLE FOR INDIRECT, INCIDENTAL, SPECIAL, CONSEQUENTIAL, OR PUNITIVE DAMAGES, INCLUDING BUT NOT LIMITED TO LOSS OF PROFITS, DATA, OR USE, RESULTING FROM THE USE OR INABILITY TO USE OUR SERVICES.
              </p>
              <p>
                Our total liability to you for any claim related to our Services shall not exceed the amount you have paid to MiCartaPro in the twelve (12) months preceding the event giving rise to the claim.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">11. Indemnification</h2>
              <p>
                You agree to indemnify, defend, and hold harmless MiCartaPro, its affiliates, directors, employees, and agents from any claim, demand, loss, liability, and expense (including attorneys' fees) arising from or related to:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Your use of the Services</li>
                <li>Your violation of these Terms</li>
                <li>Your violation of any law or right of third parties</li>
                <li>Your User Content</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">12. Termination</h2>
              <p>
                We may terminate or suspend your access to our Services immediately, without prior notice, for any reason, including but not limited to:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Violation of these Terms</li>
                <li>Fraudulent or illegal use of the Services</li>
                <li>Non-payment of applicable subscription fees</li>
                <li>For operational or business reasons</li>
              </ul>
              <p className="mt-4">
                <strong>Subscription Cancellation:</strong> You may cancel your subscription at any time by contacting us at least 48 hours before the end of your current billing period. Your cancellation will take effect at the next payment date. Upon cancellation, your access to the Services will continue until the end of your current billing period, after which your right to use the Services will cease.
              </p>
              <p className="mt-4">
                <strong>No Refunds on Unused Periods:</strong> There are no refunds on unused subscription periods. Once a subscription period has been paid, you will have access for the full period even if you cancel before the period ends.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">13. Third-Party Services</h2>
              <p>
                Our Services may integrate with third-party services, including WhatsApp. Your use of these third-party services is subject to their own terms and conditions and privacy policies. We are not responsible for the privacy practices or content of third-party services.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">14. Modifications to Terms</h2>
              <p>
                We reserve the right to modify these Terms at any time. We will notify you of material changes by posting the updated Terms on this page and updating the "Last updated" date.
              </p>
              <p>
                Your continued use of the Services after the modified Terms take effect constitutes your acceptance of the modified Terms.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">15. Governing Law and Jurisdiction</h2>
              <p>
                These Terms are governed by and construed in accordance with applicable laws, without giving effect to any principles of conflicts of law.
              </p>
              <p>
                Any dispute arising from or related to these Terms or our Services shall be resolved exclusively in the competent courts, unless otherwise agreed through arbitration.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">16. General Provisions</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">16.1. Entire Agreement</h3>
              <p>These Terms constitute the entire agreement between you and MiCartaPro regarding the use of the Services.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">16.2. Severability</h3>
              <p>If any provision of these Terms is found to be invalid or unenforceable, the remaining provisions will remain in full force and effect.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">16.3. Waiver</h3>
              <p>Our failure to exercise or enforce any right or provision of these Terms shall not constitute a waiver of such right or provision.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">16.4. Assignment</h3>
              <p>You may not assign or transfer these Terms or your rights under these Terms without our prior written consent. We may assign these Terms without restriction.</p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">17. Contact</h2>
              <p>
                If you have questions about these Terms and Conditions, you can contact us through:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56957857558</li>
                <li><strong>Email:</strong> legal@micartapro.com</li>
              </ul>
            </section>
          </div>

          <div className="mt-12 pt-8 border-t border-gray-200 flex flex-col sm:flex-row justify-between items-center gap-4">
            <Link 
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              ← Back to Home
            </Link>
            <Link 
              href="/privacy"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              View Privacy Policy →
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
