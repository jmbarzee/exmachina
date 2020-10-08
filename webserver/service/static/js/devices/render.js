// import * as THREE from '/js/devices/three.js';

var container;

var camera, scene, renderer;

var group;

var targetRotation = 0;
var targetRotationOnPointerDown = 0;

var pointerX = 0;
var pointerXOnPointerDown = 0;

var windowHalfX = window.innerWidth / 2;
init();
animate();

function init() {

    container = document.createElement( 'div' );
    document.body.appendChild( container );

    scene = new THREE.Scene();
    scene.background = new THREE.Color( 0xf0f0f0 );
    
    camera = new THREE.PerspectiveCamera( 50, window.innerWidth / window.innerHeight, 1, 1000 );
    camera.position.set( 0, 150, 500 );
    scene.add( camera );

    var light = new THREE.PointLight( 0xffffff, 0.8 );
    camera.add( light );

    group = new THREE.Group();
    scene.add( group );

    function addLineShape( shape, color, x, y, z, rx, ry, rz, s ) {
        shape.autoClose = true;

        var points = shape.getPoints();
        var geometryPoints = new THREE.BufferGeometry().setFromPoints( points );

        var line = new THREE.Line( geometryPoints, new THREE.LineBasicMaterial( { color: color } ) );
        line.position.set( x, y, z );
        line.rotation.set( rx, ry, rz );
        line.scale.set( s, s, s );
        group.add( line );
    }


    var sqLength = 80;
    var square = new THREE.Shape()
            .moveTo( 0, 0 )
            .lineTo( 0, sqLength )
            .lineTo( sqLength, sqLength )
            .lineTo( sqLength, 0 )
            .lineTo( 0, 0 );
    addLineShape(square, 0x0040f0, 0, 0, 0, 0, 0, 0, 1 )


    renderer = new THREE.WebGLRenderer( { antialias: true } );
    renderer.setPixelRatio( window.devicePixelRatio );
    renderer.setSize( window.innerWidth, window.innerHeight );
    container.appendChild( renderer.domElement );

    container.style.touchAction = 'none';
    container.addEventListener( 'pointerdown', onPointerDown, false );

    window.addEventListener( 'resize', onWindowResize, false );
}
function onWindowResize() {
    windowHalfX = window.innerWidth / 2;

    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();

    renderer.setSize( window.innerWidth, window.innerHeight );
}

function onPointerDown( event ) {
    if ( event.isPrimary === false ) return;

    pointerXOnPointerDown = event.clientX - windowHalfX;
    targetRotationOnPointerDown = targetRotation;

    document.addEventListener( 'pointermove', onPointerMove, false );
    document.addEventListener( 'pointerup', onPointerUp, false );
}

function onPointerMove( event ) {
    if ( event.isPrimary === false ) return;

    pointerX = event.clientX - windowHalfX;

    targetRotation = targetRotationOnPointerDown + ( pointerX - pointerXOnPointerDown ) * 0.02;
}

function onPointerUp() {

    if ( event.isPrimary === false ) return;

    document.removeEventListener( 'pointermove', onPointerMove );
    document.removeEventListener( 'pointerup', onPointerUp );
}

function animate() {
    requestAnimationFrame( animate );

    render();
}

function render() {
    group.rotation.y += ( targetRotation - group.rotation.y ) * 0.05;
    renderer.render( scene, camera );
}
